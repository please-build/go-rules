package goget

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

var client = http.DefaultClient

type moduleVersion struct {
	mod, ver string
}

type getter struct {
	queryResults map[moduleVersion]*modfile.File
	proxyUrl     string
}

func newGetter() *getter {
	return &getter{
		queryResults: map[moduleVersion]*modfile.File{},
		proxyUrl:     "https://proxy.golang.org",
	}
}

func (g *getter) getGoMod(mod, ver string) (*modfile.File, error) {
	modVer := moduleVersion{mod, ver}
	if modFile, ok := g.queryResults[modVer]; ok {
		return modFile, nil
	}

	file := fmt.Sprintf("%s/%s/@v/%s.mod", g.proxyUrl, mod, ver)
	resp, err := client.Get(file)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%v %v: \n%v", file, resp.StatusCode, string(body))
	}

	modFile, err := modfile.Parse(file, body, nil)
	if err != nil {
		return nil, err
	}

	g.queryResults[modVer] = modFile
	return modFile, nil
}

// getGoModWithFallback attempts to get a go.mod for the given module and
// version with fallback for supporting modules with case insensitivity.
func (g *getter) getGoModWithFallback(mod, version string) (*modfile.File, error) {
	modVersionsToAttempt := map[string]string{
		mod: version,
	}

	// attempt lowercasing entire mod string for packages like:
	// - `github.com/Sirupsen/logrus` -> `github.com/sirupsen/logrus`.
	// https://github.com/sirupsen/logrus/issues/543
	modVersionsToAttempt[strings.ToLower(mod)] = version

	var errs error
	for mod, version := range modVersionsToAttempt {
		modFile, err := g.getGoMod(mod, version)
		if err != nil {
			// TODO: when we upgrade to Go 1.20, use `errors.Join(...)`.
			errs = fmt.Errorf("%w: %w", errs, err)
			continue
		}

		return modFile, nil
	}

	return nil, errs
}

func (g *getter) getDeps(deps map[string]string, mod, version string) error {
	modFile, err := g.getGoModWithFallback(mod, version)
	if err != nil {
		return err
	}

	for _, req := range modFile.Require {
		oldVer, ok := deps[req.Mod.Path]
		if !ok || semver.Compare(oldVer, req.Mod.Version) < 0 {
			deps[req.Mod.Path] = req.Mod.Version
			if err := g.getDeps(deps, req.Mod.Path, req.Mod.Version); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *getter) goGet(mods []string) error {
	deps := map[string]string{}
	for _, mod := range mods {
		parts := strings.Split(mod, "@")
		deps[parts[0]] = parts[1]
	}

	for mod, ver := range deps {
		if err := g.getDeps(deps, mod, ver); err != nil {
			return err
		}
	}

	for mod, ver := range deps {
		fmt.Printf("go_repo(module=\"%s\", version=\"%s\")\n", mod, ver)
	}
	return nil
}

// GoGet is used to spit out a new go_get rule. The plan is to build this out into a tool to add new third party
// modules to the repo.
func GoGet(mods []string) error {
	return newGetter().goGet(mods)
}

func GetMod(path string) error {
	g := newGetter()
	bs, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	modFile, err := modfile.ParseLax(path, bs, nil)
	if err != nil {
		return err
	}

	paths := make([]string, len(modFile.Require))
	for i, req := range modFile.Require {
		paths[i] = fmt.Sprintf("%v@%v", req.Mod.Path, req.Mod.Version)
	}

	return g.goGet(paths)
}

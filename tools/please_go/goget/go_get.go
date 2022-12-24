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

	file := fmt.Sprintf("%s/%s/@v/%s.mod", g.proxyUrl, strings.ToLower(mod), ver)
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

func (g *getter) getDeps(deps map[string]string, mod, version string) error {
	modFile, err := g.getGoMod(mod, version)
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
		fmt.Printf("go_get(module=\"%s\", version=\"%s\")\n", mod, ver)
	}
	return nil
}

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

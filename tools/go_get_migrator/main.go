package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/buildtools/build"
	"github.com/please-build/go-rules/tools/please_go/generate"
	"golang.org/x/mod/semver"
)

func getModules() map[string]string {
	cmd := exec.Command("plz", "query", "print", "--label=go_module:", "//...")
	var stdErr bytes.Buffer
	var stdOut bytes.Buffer

	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	err := cmd.Run()
	if err != nil {
		panic(fmt.Errorf("failed to get moduels: %v, StdErr: \n%v", err, stdErr.String()))
	}

	reqPairs := strings.Split(strings.TrimSpace(stdOut.String()), "\n")
	reqs := map[string]string{}

	for _, pair := range reqPairs {
		parts := strings.Split(pair, "@")
		if len(parts) != 2 {
			panic(fmt.Errorf("invalid mod requirement: %v", pair))
		}

		if ver, ok := reqs[parts[0]]; !ok || semver.Compare(ver, parts[1]) < 0 {
			reqs[parts[0]] = parts[1]
		}

	}
	return reqs
}

func parseBuildFile(file string) (*build.File, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return build.ParseBuild(file, data)
}

func findRepoRoot() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return findRepoRootInPath(path)
}

func findRepoRootInPath(path string) (string, error) {
	if path == "." {
		return "", errors.New("failed to find repo root")
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.Name() == ".plzconfig" {
			return path, nil
		}
	}

	return findRepoRootInPath(filepath.Dir(path))
}

func getImportPath() (string, error) {
	out, err := exec.Command("plz", "query", "config", "plugin.go.importpath").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running command, \"plz query config pugin.go.importpath\": %v\n%v", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

type TargetJSON struct {
	Labels []string `json:"labels"`
}

func getGoPackageTargetMapping() (map[string]string, error) {
	out, err := exec.Command("plz", "query", "print", "--field", "labels", "--json", "-i", "go", "//...").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to query go targets: %v\n%v", err, string(out))
	}

	targets := map[string]*TargetJSON{}
	err = json.Unmarshal(out, &targets)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string, len(targets))
	for label, target := range targets {
		for _, l := range target.Labels {
			if strings.HasPrefix(l, "go_package:") {
				ret[label] = strings.TrimPrefix(l, "go_package:")
			}
		}
	}
	return ret, nil
}

var opts = struct {
	Usage string

	Update struct {
		Args struct {
			Packages []string `positional-arg-name:"packages" description:"The packages to compile"`
		} `positional-args:"true" required:"true"`
	} `command:"update" alias:"i" description:"Updates some build files"`
	Mod struct {
		Sync struct {
		} `command:"sync" alias:"i" description:"synchronises Please with a go.mod file"`
	} `command:"mod" alias:"i" description:""`
}{
	Usage: `
`,
}

func main() {
	// TODO(jpoole): configure the third party build file path
	file, err := parseBuildFile("third_party/go/BUILD")
	if err != nil {
		panic(err)
	}

	var stmts []build.Expr
	for _, stmt := range file.Stmt {
		if expr, ok := stmt.(*build.CallExpr); ok {
			if ident, ok := expr.X.(*build.Ident); ok && ident.Name == "go_module" || ident.Name == "go_mod_download" || ident.Name == "go_get" {
				// Strip out all the existing rules
				continue
			}
		}
		stmts = append(stmts, stmt)
	}

	file.Stmt = stmts

	modReqs := getModules()

	// TODO(jpoole): have some ordering so this is deterministic
	modules := make([]string, 0, len(modReqs))
	for mod, ver := range modReqs {
		rule := &Target{Kind: "go_get", Attributes: map[string]build.Expr{}}
		rule.Attributes["module"] = Expr(mod)
		rule.Attributes["version"] = Expr(ver)

		file.Stmt = append(file.Stmt, rule.ToCallExpr())
		modules = append(modules, mod)
	}

	if err := os.WriteFile("third_party/go/BUILD", build.Format(file), 660); err != nil {
		panic(err)
	}

	reporoot, err := findRepoRoot()
	if err != nil {
		panic(err)
	}

	importPath, err := getImportPath()
	if err != nil {
		panic(err)
	}

	g := generate.New(reporoot, "third_party/go", []string{"BUILD", "BUILD.plz"}, modules)

	if err := g.Update(importPath, os.Args[1:]); err != nil {
		panic(err)
	}
}

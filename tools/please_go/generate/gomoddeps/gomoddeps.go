// Package gomoddeps parses dependencies and replacements from the host and module go.mod files.
package gomoddeps

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"golang.org/x/mod/modfile"
)

// GetCombinedDepsAndReplacements returns dependencies and replacements after inspecting both
// the host and the module go.mod files.
// Module's replacement are only returned if there is no host go.mod file.
func GetCombinedDepsAndReplacements(hostGoModPath, moduleGoModPath string) ([]string, map[string]string, error) {
	var err error

	hostDeps := []string{}
	hostReplacements := map[string]string{}
	if hostGoModPath != "" {
		hostDeps, hostReplacements, err = getDepsAndReplacements(hostGoModPath, false)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read host repo go.mod %q: %w", hostGoModPath, err)
		}
	}

	var moduleDeps []string
	var moduleReplacements = map[string]string{}
	useLaxParsingForModule := true
	if hostGoModPath == "" {
		// If we're only considering the module then we want to extract the replacement's as well (lax mode
		// doesn't parse them).
		useLaxParsingForModule = false
	}
	moduleDeps, moduleReplacements, err = getDepsAndReplacements(moduleGoModPath, useLaxParsingForModule)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return hostDeps, hostReplacements, nil
		}
		return nil, nil, fmt.Errorf("failed to read module go.mod %q: %w", moduleGoModPath, err)
	}

	var replacements map[string]string
	if hostGoModPath == "" {
		replacements = moduleReplacements
	} else {
		replacements = hostReplacements
	}

	return append(hostDeps, moduleDeps...), replacements, nil
}

// getDepsAndReplacements parses the go.mod file and returns all the dependencies
// and replacement directives from it.
func getDepsAndReplacements(goModPath string, useLaxParsing bool) ([]string, map[string]string, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, nil, err
	}

	var modFile *modfile.File
	if useLaxParsing {
		modFile, err = modfile.ParseLax(goModPath, data, nil)
	} else {
		modFile, err = modfile.Parse(goModPath, data, nil)
	}
	if err != nil {
		return nil, nil, err
	}

	moduleDeps := make([]string, 0, len(modFile.Require))
	// TODO we could probably validate these are known modules
	for _, req := range modFile.Require {
		moduleDeps = append(moduleDeps, req.Mod.Path)
	}

	replacements := make(map[string]string, len(modFile.Replace))
	for _, replace := range modFile.Replace {
		replacements[replace.Old.Path] = replace.New.Path
	}

	return moduleDeps, replacements, nil
}

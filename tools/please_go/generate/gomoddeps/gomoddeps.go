package gomoddeps

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"golang.org/x/mod/modfile"
)

// GetCombinedDepsAndRequirements returns dependencies and replacements after inspecting both
// the host and the module go.mod files.
// Module's replacement are only returned if there is no host go.mod file.
func GetCombinedDepsAndRequirements(hostGoModPath, moduleGoModPath string) ([]string, map[string]string, error) {
	var err error

	hostDeps := []string{}
	hostReplacements := map[string]string{}
	if hostGoModPath != "" {
		hostDeps, hostReplacements, err = getDepsAndReplacements(hostGoModPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read host repo go.mod %q: %w", hostGoModPath, err)
		}
	}

	var moduleDeps []string
	var moduleReplacements = map[string]string{}
	moduleDeps, moduleReplacements, err = getDepsAndReplacements(moduleGoModPath)
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
func getDepsAndReplacements(goModPath string) ([]string, map[string]string, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, nil, err
	}
	modFile, err := modfile.Parse(goModPath, data, nil)
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

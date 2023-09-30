package generate

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	bazelbuild "github.com/bazelbuild/buildtools/build"
	bazeledit "github.com/bazelbuild/buildtools/edit"
)

type Generate struct {
	moduleName         string
	moduleArg          string
	srcRoot            string
	buildContext       build.Context
	modFile            string
	buildFileNames     []string
	moduleDeps         []string
	pluginTarget       string
	replace            map[string]string
	knownImportTargets map[string]string // cache these so we don't end up looping over all the modules for every import
	thirdPartyFolder   string
	install            []string
}

func New(srcRoot, thirdPartyFolder, modFile, module, version string, buildFileNames, moduleDeps, install []string) *Generate {
	moduleArg := module
	if version != "" {
		moduleArg += "@" + version
	}
	return &Generate{
		srcRoot:            srcRoot,
		buildContext:       build.Default,
		buildFileNames:     buildFileNames,
		moduleDeps:         moduleDeps,
		modFile:            modFile,
		knownImportTargets: map[string]string{},
		thirdPartyFolder:   thirdPartyFolder,
		install:            install,
		moduleName:         module,
		moduleArg:          moduleArg,
	}
}

// Generate generates a new Please project at the src root. It will walk through the directory tree generating new BUILD
// files. This is primarily intended to generate a please subrepo for third party code.
func (g *Generate) Generate() error {
	if err := g.readGoMod(g.modFile); err != nil {
		return err
	}
	if err := g.writeConfig(); err != nil {
		return err
	}
	if err := g.generateAll(g.srcRoot); err != nil {
		return err
	}
	return g.writeInstallFilegroup()
}

func (g *Generate) installTargets() []string {
	var targets []string

	for _, i := range g.install {
		dir := filepath.Join(g.srcRoot, i)
		if strings.HasSuffix(dir, "/...") {
			targets = append(targets, g.targetsInDir(strings.TrimSuffix(dir, "/..."))...)
		} else {
			targets = append(targets, g.libTargetForPleasePackage(i))
		}
	}
	return targets
}

func (g *Generate) targetsInDir(dir string) []string {
	var ret []string
	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		// The assumption here is that if we generated a BUILD file, then we would have generated a go_library() target
		// for that package. Currently, we don't generate BUILD files for any other reason so this assumption holds
		// true. We may want to check that the BUILD file contains a go_library() target otherwise.
		if g.isBuildFile(path) {
			ret = append(ret, g.libTargetForPleasePackage(trimPath(filepath.Dir(path), g.srcRoot)))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return ret
}

func (g *Generate) isBuildFile(file string) bool {
	base := filepath.Base(file)
	for _, file := range g.buildFileNames {
		if base == file {
			return true
		}
	}
	return false
}

func (g *Generate) writeInstallFilegroup() error {
	buildFile, err := parseOrCreateBuildFile(g.srcRoot, g.buildFileNames)
	if err != nil {
		return err
	}

	rule := NewRule("filegroup", "installs")
	rule.SetAttr("exported_deps", NewStringList(g.installTargets()))
	rule.SetAttr("visibility", NewStringList([]string{"PUBLIC"}))

	buildFile.Stmt = append(buildFile.Stmt, rule.Call)

	return saveBuildFile(buildFile)
}

// readGoMod reads the module dependencies
func (g *Generate) readGoMod(path string) error {
	if path == "" {
		path = filepath.Join(g.srcRoot, "go.mod")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	modFile, err := modfile.ParseLax(path, data, nil)
	if err != nil {
		return err
	}

	// TODO we could probably validate these are known modules
	for _, req := range modFile.Require {
		g.moduleDeps = append(g.moduleDeps, req.Mod.Path)
	}

	g.moduleDeps = append(g.moduleDeps, g.moduleName)

	g.replace = make(map[string]string, len(modFile.Replace))
	for _, replace := range modFile.Replace {
		g.replace[replace.Old.Path] = replace.New.Path
	}
	return nil
}

func (g *Generate) writeConfig() error {
	file, err := os.Create(filepath.Join(g.srcRoot, ".plzconfig"))
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "[Plugin \"go\"]")
	fmt.Fprintln(file, "Target=@//plugins:go")
	fmt.Fprintf(file, "ImportPath=%s\n", g.moduleName)
	return nil
}

func (g *Generate) generateAll(dir string) error {
	return filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "testdata" {
				return filepath.SkipDir
			}
			if path != dir && strings.HasPrefix(info.Name(), "_") {
				return filepath.SkipDir
			}

			if err := g.generate(trimPath(path, g.srcRoot)); err != nil {
				switch err.(type) {
				case *build.NoGoError:
					// We might walk into a dir that has no .go files for the current arch. This shouldn't
					// be an error so we just eat this
					return nil
				default:
					return err
				}
			}
		}
		return nil
	})
}

func (g *Generate) pkgDir(target string) string {
	p := strings.TrimPrefix(target, g.moduleName)
	return filepath.Join(g.srcRoot, p)
}

func (g *Generate) importDir(target string) (*build.Package, error) {
	dir := filepath.Join(os.Getenv("TMP_DIR"), g.pkgDir(target))
	return g.buildContext.ImportDir(dir, 0)
}

func (g *Generate) generate(dir string) error {
	pkg, err := g.importDir(dir)
	if err != nil {
		return err
	}

	lib := g.libRule(pkg, dir)
	if lib == nil {
		return nil
	}

	return g.createBuildFile(dir, lib)
}

func (g *Generate) matchesInstall(dir string) bool {
	for _, i := range g.install {
		i := filepath.Join(g.srcRoot, i)
		pkgDir := g.pkgDir(dir)

		if strings.HasSuffix(i, "/...") {
			i = strings.TrimSuffix(i, "/...")
			return strings.HasPrefix(pkgDir, i)
		}
		return i == pkgDir
	}
	return false
}

func (g *Generate) rule(rule *Rule) *bazelbuild.Rule {
	r := NewRule(rule.kind, rule.name)
	populateRule(r, rule)
	r.SetAttr("visibility", NewStringList([]string{"PUBLIC"}))

	return r
}

// parseOrCreateBuildFile loops through the available build file names to create a new build file or open the existing
// one.
func parseOrCreateBuildFile(path string, fileNames []string) (*bazelbuild.File, error) {
	for _, name := range fileNames {
		filePath := filepath.Join(path, name)
		if f, err := os.Lstat(filePath); os.IsNotExist(err) {
			return bazelbuild.ParseBuild(filePath, nil)
		} else if !f.IsDir() {
			bs, err := os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			return bazelbuild.ParseBuild(filePath, bs)
		}
	}
	return nil, fmt.Errorf("folders exist with the build file names in directory %v %v", path, fileNames)
}

func saveBuildFile(buildFile *bazelbuild.File) error {
	f, err := os.Create(buildFile.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(bazelbuild.Format(buildFile))
	return err
}

func (g *Generate) createBuildFile(pkg string, rule *Rule) error {
	buildFile, err := parseOrCreateBuildFile(g.pkgDir(pkg), g.buildFileNames)
	if err != nil {
		return err
	}

	var subincludes []bazelbuild.Expr
	if strings.HasPrefix(rule.kind, "cgo") {
		subincludes = []bazelbuild.Expr{NewStringExpr("///go//build_defs:cgo")}
	} else {
		subincludes = []bazelbuild.Expr{NewStringExpr("///go//build_defs:go")}
	}

	buildFile.Stmt = []bazelbuild.Expr{
		&bazelbuild.CallExpr{
			X:    &bazelbuild.Ident{Name: "subinclude"},
			List: subincludes,
		},
	}

	buildFile.Stmt = append(buildFile.Stmt, g.rule(rule).Call)

	return saveBuildFile(buildFile)
}

func NewRule(kind, name string) *bazelbuild.Rule {
	rule, _ := bazeledit.ExprToRule(&bazelbuild.CallExpr{
		X:    &bazelbuild.Ident{Name: kind},
		List: []bazelbuild.Expr{},
	}, kind)

	rule.SetAttr("name", NewStringExpr(name))

	return rule
}

func NewStringExpr(s string) *bazelbuild.StringExpr {
	return &bazelbuild.StringExpr{Value: s}
}

func NewStringList(ss []string) *bazelbuild.ListExpr {
	l := new(bazelbuild.ListExpr)
	for _, s := range ss {
		l.List = append(l.List, NewStringExpr(s))
	}
	return l
}

func packageKind(pkg *build.Package) string {
	cgo := len(pkg.CgoFiles) > 0
	if pkg.IsCommand() && cgo {
		return "cgo_binary"
	}
	if pkg.IsCommand() {
		return "go_binary"
	}
	if cgo {
		return "cgo_library"
	}
	return "go_library"
}

func (g *Generate) depTargets(imports []string) []string {
	deps := make([]string, 0)
	for _, path := range imports {
		target := g.depTarget(path)
		if target == "" {
			continue
		}
		deps = append(deps, target)
	}
	return deps
}

func (g *Generate) libRule(pkg *build.Package, dir string) *Rule {
	if len(pkg.GoFiles) == 0 && len(pkg.CgoFiles) == 0 {
		return nil
	}

	name := nameForLibInPkg(g.moduleName, trimPath(dir, g.srcRoot))

	return &Rule{
		name:          name,
		kind:          packageKind(pkg),
		srcs:          pkg.GoFiles,
		module:        g.moduleArg,
		cgoSrcs:       pkg.CgoFiles,
		compilerFlags: pkg.CgoCFLAGS,
		linkerFlags:   pkg.CgoLDFLAGS,
		pkgConfigs:    pkg.CgoPkgConfig,
		asmFiles:      pkg.SFiles,
		hdrs:          pkg.HFiles,
		deps:          g.depTargets(pkg.Imports),
		embedPatterns: pkg.EmbedPatterns,
	}
}

func (g *Generate) depTarget(importPath string) string {
	if target, ok := g.knownImportTargets[importPath]; ok {
		return target
	}

	if replacement, ok := g.replace[importPath]; ok {
		target := g.depTarget(replacement)
		g.knownImportTargets[importPath] = target
		return target
	}

	module := ""
	for _, mod := range append(g.moduleDeps, g.moduleName) {
		if strings.HasPrefix(importPath, mod) {
			if len(module) < len(mod) {
				module = mod
			}
		}
	}

	if module == "" {
		// If we can't find this import, we can return nothing and the build rule will fail at build time reporting a
		// sensible error. It may also be an import from the go SDK which is fine.
		return ""
	}

	subrepoName := g.subrepoName(module)
	packageName := trimPath(importPath, module)
	name := nameForLibInPkg(module, packageName)

	target := buildTarget(name, packageName, subrepoName)
	g.knownImportTargets[importPath] = target
	return target
}

// nameForLibInPkg returns the lib target name for a target in pkg. The pkg should be the relative pkg part excluding
// the module, e.g. pkg would be asset, and module would be github.com/stretchr/testify for
// github.com/stretchr/testify/assert,
func nameForLibInPkg(module, pkg string) string {
	name := filepath.Base(pkg)
	if pkg == "" || pkg == "." {
		name = filepath.Base(module)
	}

	if name == "all" {
		return "lib"
	}

	return name
}

// trimPath is like strings.TrimPrefix but is path aware. It removes base from target if target starts with base,
// otherwise returns target unmodified.
func trimPath(target, base string) string {
	baseParts := strings.Split(filepath.Clean(base), "/")
	targetParts := strings.Split(filepath.Clean(target), "/")

	if len(targetParts) < len(baseParts) {
		return target
	}

	for i := range baseParts {
		if baseParts[i] != targetParts[i] {
			return target
		}
	}
	return strings.Join(targetParts[len(baseParts):], "/")
}

// libTargetForPleasePackage returns the build label for the go_library() target that would be generated for a package
// at this path within the generated Please repo.
func (g *Generate) libTargetForPleasePackage(pkg string) string {
	return buildTarget(nameForLibInPkg(g.moduleName, pkg), pkg, "")
}

func (g *Generate) subrepoName(module string) string {
	if g.moduleName == module {
		return ""
	}
	return filepath.Join(g.thirdPartyFolder, strings.ReplaceAll(module, "/", "_"))
}

func buildTarget(name, pkgDir, subrepo string) string {
	bs := new(strings.Builder)
	if subrepo != "" {
		bs.WriteString("///")
		bs.WriteString(subrepo)
	}

	// Bit of a special case here where we assume all build targets are absolute which is fine for our use case.
	bs.WriteString("//")

	if pkgDir == "." {
		pkgDir = ""
	}

	if pkgDir != "" {
		bs.WriteString(pkgDir)
		if filepath.Base(pkgDir) != name {
			bs.WriteString(":")
			bs.WriteString(name)
		}
	} else {
		bs.WriteString(":")
		bs.WriteString(name)
	}
	return bs.String()
}

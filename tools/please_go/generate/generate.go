package generate

import (
	"fmt"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	bazelbuild "github.com/bazelbuild/buildtools/build"
	bazeledit "github.com/bazelbuild/buildtools/edit"
)

type Generate struct {
	moduleName         string
	srcRoot            string
	buildContext       build.Context
	buildFileNames     []string
	moduleDeps         []string
	pluginTarget       string
	replace            map[string]string
	knownImportTargets map[string]string // cache these so we don't end up looping over all the modules for every import
	thirdPartyFolder   string
	install            []string
}

func New(srcRoot, thirdPartyFolder string, buildFileNames, moduleDeps, install []string) *Generate {
	return &Generate{
		srcRoot:            srcRoot,
		buildContext:       build.Default,
		buildFileNames:     buildFileNames,
		moduleDeps:         moduleDeps,
		knownImportTargets: map[string]string{},
		thirdPartyFolder:   thirdPartyFolder,
		install:            install,
	}
}

// Generate generates a new Please project at the src root. It will walk through the directory tree generating new BUILD
// files. This is primarily intended to generate a please subrepo for third party code.
func (g *Generate) Generate() error {
	if err := g.readGoMod(); err != nil {
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
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		// The assumption here is that if we generated a BUILD file, then we would have generated a go_library() target
		// for that package. Currently, we don't generate BUILD files for any other reason so this assumption holds
		// true. We may want to check that the BUILD file contains a go_library() target otherwise.
		if g.isBuildFile(path) {
			pleasePkgDir := strings.TrimPrefix(strings.TrimPrefix(filepath.Dir(path), g.srcRoot), "/")
			ret = append(ret, g.libTargetForPleasePackage(pleasePkgDir))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return ret
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
func (g *Generate) readGoMod() error {
	path := filepath.Join(g.srcRoot, "go.mod")
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

	g.moduleName = modFile.Module.Mod.Path
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
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "testdata" {
				return filepath.SkipDir
			}
			if err := g.generate(filepath.Clean(strings.TrimPrefix(path, g.srcRoot))); err != nil {
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

	lib := g.libRule(pkg)
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

func (g *Generate) libRule(pkg *build.Package) *Rule {
	// The name of the target should match the dir it's in, or the basename of the module if it's in the repo root.
	name := filepath.Base(pkg.Dir)
	if strings.HasSuffix(pkg.Dir, g.srcRoot) || name == "" {
		name = filepath.Base(g.moduleName)
	}

	if len(pkg.GoFiles) == 0 && len(pkg.CgoFiles) == 0 {
		return nil
	}

	return &Rule{
		name:          name,
		kind:          packageKind(pkg),
		srcs:          pkg.GoFiles,
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

func (g *Generate) testRule(pkg *build.Package, prodRule *Rule) *Rule {
	// The name of the target should match the dir it's in, or the basename of the module if it's in the repo root.
	name := filepath.Base(pkg.Dir)
	if strings.HasSuffix(pkg.Dir, g.srcRoot) || name == "" {
		name = filepath.Base(g.moduleName)
	}

	if len(pkg.TestGoFiles) == 0 {
		return nil
	}
	deps := g.depTargets(pkg.TestImports)

	if prodRule != nil {
		deps = append(deps, ":"+prodRule.name)
	}
	// TODO(jpoole): handle external tests
	return &Rule{
		name:          name,
		kind:          "go_test",
		srcs:          pkg.TestGoFiles,
		deps:          deps,
		embedPatterns: pkg.TestEmbedPatterns,
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
	packageName := strings.TrimPrefix(importPath, module)
	packageName = strings.TrimPrefix(packageName, "/")
	name := filepath.Base(packageName)
	if packageName == "" {
		name = filepath.Base(module)
	}

	target := buildTarget(name, packageName, subrepoName)
	g.knownImportTargets[importPath] = target
	return target
}

// libTargetForPleasePackage returns the build label for the go_library() target that would be generated for a package
// at this path withing the generated Please repo.
func (g *Generate) libTargetForPleasePackage(pkg string) string {
	if pkg == "" || pkg == "." {
		return buildTarget(filepath.Base(g.moduleName), "", "")
	}
	return buildTarget(filepath.Base(pkg), pkg, "")
}

func (g *Generate) subrepoName(module string) string {
	if g.moduleName == module {
		return ""
	}
	return filepath.Join(g.thirdPartyFolder, strings.ReplaceAll(module, "/", "_"))
}

func buildTarget(name, pkg, subrepo string) string {
	if name == "" {
		name = filepath.Base(pkg)
	}
	target := fmt.Sprintf("%v:%v", pkg, name)
	if subrepo == "" {
		return fmt.Sprintf("//%v", target)
	}
	return fmt.Sprintf("///%v//%v", subrepo, target)
}

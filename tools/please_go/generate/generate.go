package generate

import (
	"bufio"
	"fmt"
	"go/build"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	bazelbuild "github.com/bazelbuild/buildtools/build"
	bazeledit "github.com/bazelbuild/buildtools/edit"

	"github.com/please-build/go-rules/tools/please_go/generate/gomoddeps"
)

type Generate struct {
	moduleName         string
	moduleArg          string
	srcRoot            string
	subrepo            string
	buildContext       build.Context
	hostModFile        string
	buildFileNames     []string
	moduleDeps         []string
	replace            map[string]string
	knownImportTargets map[string]string // cache these so we don't end up looping over all the modules for every import
	thirdPartyFolder   string
	install            []string
	labels             []string
}

func New(srcRoot, thirdPartyFolder, hostModFile, module, version, subrepo string, buildFileNames, moduleDeps, install []string, buildTags []string, labels []string) *Generate {
	moduleArg := module
	if version != "" {
		moduleArg += "@" + version
	}

	ctxt := build.Default
	ctxt.BuildTags = buildTags

	return &Generate{
		srcRoot:            srcRoot,
		buildContext:       ctxt,
		buildFileNames:     buildFileNames,
		moduleDeps:         moduleDeps,
		hostModFile:        hostModFile,
		knownImportTargets: map[string]string{},
		thirdPartyFolder:   thirdPartyFolder,
		install:            install,
		moduleName:         module,
		moduleArg:          moduleArg,
		subrepo:            subrepo,
		labels:             labels,
	}
}

// Generate generates a new Please project at the src root. It will walk through the directory tree generating new BUILD
// files. This is primarily intended to generate a please subrepo for third party code.
func (g *Generate) Generate() error {
	deps, replacements, err := gomoddeps.GetCombinedDepsAndReplacements(g.hostModFile, path.Join(g.srcRoot, "go.mod"))
	if err != nil {
		return err
	}
	// It's important to not override g.moduleDeps as it can already contains dependencies configured
	// when `Generate` was constructed.
	g.moduleDeps = append(g.moduleDeps, deps...)
	g.moduleDeps = append(g.moduleDeps, g.moduleName)
	g.replace = replacements

	if err := g.writeConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	if err := g.parseImportConfigs(); err != nil {
		return fmt.Errorf("failed to parse import configs: %w", err)
	}

	if err := g.generateAll(g.srcRoot); err != nil {
		return fmt.Errorf("failed to generate BUILD files: %w", err)
	}
	return g.writeInstallFilegroup()
}

// parseImportConfigs walks through the build dir looking for .importconfig files, parsing the # please:target //foo:bar
// comments to generate the known imports. These are the deps that are passed to the go_repo e.g. for legacy go_module
// rules.
func (g *Generate) parseImportConfigs() error {
	return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".importconfig" {
			target, pkgs, err := parseImportConfig(path)
			if err != nil {
				return err
			}
			if target == "" {
				return nil
			}
			for _, p := range pkgs {
				g.knownImportTargets[p] = target
			}
		}
		return nil
	})
}

func parseImportConfig(path string) (string, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	target := ""
	var imports []string

	importCfg := bufio.NewScanner(f)
	for importCfg.Scan() {
		line := importCfg.Text()
		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "# please:target ") {
				target = strings.TrimSpace(strings.TrimPrefix(line, "# please:target "))
				if !strings.HasPrefix(target, "///") {
					target = "@" + target
				}
			}
			continue
		}
		parts := strings.Split(strings.TrimPrefix(line, "packagefile "), "=")
		imports = append(imports, parts[0])
	}
	return target, imports, nil
}

func (g *Generate) installTargets() ([]string, error) {
	var targets []string

	for _, i := range g.install {
		dir := filepath.Join(g.srcRoot, i)
		if strings.HasSuffix(dir, "/...") {
			ts, err := g.targetsInDir(strings.TrimSuffix(dir, "/..."))
			if err != nil {
				return nil, err
			}
			targets = append(targets, ts...)
		} else {
			t, err := g.libTargetForBuildPackage(i)
			if err != nil {
				return nil, err
			}
			if t == "" {
				return nil, fmt.Errorf("couldn't find install package %v", i)
			}
			targets = append(targets, t)
		}
	}
	return targets, nil
}

func (g *Generate) targetsInDir(dir string) ([]string, error) {
	var ret []string
	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if g.isBuildFile(path) {
			t, err := g.libTargetForBuildFile(trimPath(path, g.srcRoot))
			if err != nil {
				return err
			}
			if t != "" {
				ret = append(ret, t)
			}
		}
		return nil
	})
	return ret, err
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
	installTargets, err := g.installTargets()
	if err != nil {
		return fmt.Errorf("failed to generate install targets: %v", err)
	}
	rule.SetAttr("exported_deps", NewStringList(installTargets))
	rule.SetAttr("visibility", NewStringList([]string{"PUBLIC"}))

	buildFile.Stmt = append(buildFile.Stmt, rule.Call)

	return saveBuildFile(buildFile)
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
	for _, t := range g.buildContext.BuildTags {
		fmt.Fprintf(file, "BuildTags=%s\n", t)
	}
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
	pkg, err := g.buildContext.ImportDir(dir, 0)
	if err != nil {
		return nil, err
	}
	// We also need to discover & attach any .a files in the directory; some libraries use these
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	pkg.IgnoredOtherFiles = nil
	for _, entry := range entries {
		if name := entry.Name(); strings.HasSuffix(name, ".a") {
			pkg.IgnoredOtherFiles = append(pkg.IgnoredOtherFiles, name)
		}
	}
	return pkg, nil
}

func (g *Generate) generate(dir string) error {
	pkg, err := g.importDir(dir)
	if err != nil {
		return err
	}

	// filter out pkg.GoFiles based on build tags
	var goFiles []string
	for _, f := range pkg.GoFiles {
		match, err := g.buildContext.MatchFile(pkg.Dir, f)
		if err != nil {
			return err
		}
		if match {
			goFiles = append(goFiles, f)
		}
	}

	pkg.GoFiles = goFiles
	lib := g.ruleForPackage(pkg, dir)
	if lib == nil {
		return nil
	}

	return g.createBuildFile(dir, lib, pkg.IgnoredOtherFiles)
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
	r.SetAttr("labels", NewStringList(g.labels))
	if rule.kind == "go_library" {
		r.SetAttr("cover", &bazelbuild.Ident{Name: "False"})
	}

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

func (g *Generate) createBuildFile(pkg string, rule *Rule, aFiles []string) error {
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

	if len(aFiles) != 0 {
		filegroup := NewRule("filegroup", "a_files")
		filegroup.SetAttr("srcs", NewStringList(aFiles))
		buildFile.Stmt = append(buildFile.Stmt, filegroup.Call)
	}

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

func (g *Generate) ruleForPackage(pkg *build.Package, dir string) *Rule {
	if len(pkg.GoFiles) == 0 && len(pkg.CgoFiles) == 0 {
		return nil
	}

	name := nameForLibInPkg(g.moduleName, trimPath(dir, g.srcRoot))
	deps := g.depTargets(pkg.Imports)
	if len(pkg.IgnoredOtherFiles) != 0 {
		deps = append(deps, ":a_files")
	}

	return &Rule{
		name:          name,
		kind:          packageKind(pkg),
		srcs:          pkg.GoFiles,
		module:        g.moduleArg,
		subrepo:       g.subrepo,
		cgoSrcs:       pkg.CgoFiles,
		cSrcs:         pkg.CFiles,
		compilerFlags: pkg.CgoCFLAGS,
		linkerFlags:   orderLinkerFlags(pkg.CgoLDFLAGS),
		pkgConfigs:    pkg.CgoPkgConfig,
		asmFiles:      pkg.SFiles,
		hdrs:          pkg.HFiles,
		deps:          deps,
		embedPatterns: pkg.EmbedPatterns,
		isCMD:         pkg.IsCommand(),
	}
}

// orderLinkerFlags collapses linker flags into one to enforce a consistent ordering
func orderLinkerFlags(in []string) []string {
	if len(in) > 0 {
		return []string{strings.Join(in, " ")}
	}
	return nil
}

func (g *Generate) depTarget(importPath string) string {
	if target, ok := g.knownImportTargets[importPath]; ok {
		return target
	}

	if replacement, ok := g.replace[importPath]; ok && replacement != importPath {
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

// libTargetForBuildFile finds the go_library or cgo_library target in the package
func (g *Generate) libTargetForBuildFile(path string) (string, error) {
	bs, err := os.ReadFile(filepath.Join(g.srcRoot, path))
	if err != nil {
		return "", err
	}
	file, err := bazelbuild.ParseBuild(path, bs)
	if err != nil {
		return "", err
	}

	libs := append(file.Rules("go_library"), file.Rules("cgo_library")...)
	if len(libs) >= 1 {
		if len(libs) != 1 {
			log.Fatalf("more than one go library in installed package %v", path)
		}
		return buildTarget(libs[0].Name(), filepath.Dir(path), ""), nil
	}
	return "", nil
}

func (g *Generate) subrepoName(module string) string {
	if g.moduleName == module {
		return ""
	}
	return filepath.Join(g.thirdPartyFolder, strings.ReplaceAll(module, "/", "_"))
}

func (g *Generate) libTargetForBuildPackage(i string) (string, error) {
	entries, err := os.ReadDir(filepath.Join(g.srcRoot, i))
	if err != nil {
		return "", err
	}

	for _, e := range entries {
		if g.isBuildFile(e.Name()) {
			t, err := g.libTargetForBuildFile(filepath.Join(i, e.Name()))
			if err != nil {
				return "", err
			}
			return t, nil
		}
	}
	return "", nil
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

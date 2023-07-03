package test

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

type testDescr struct {
	Package        string
	Main           string
	TestFunctions  []string
	BenchFunctions []string
	FuzzFunctions  []string
	Examples       []*doc.Example
	CoverVars      []CoverVar
	Imports        []string
	Coverage       bool
	Benchmark      bool
	HasFuzz        bool
}

// WriteTestMain templates a test main file from the given sources to the given output file.
func WriteTestMain(testPackage string, sources []string, output string, coverage bool, coverVars []CoverVar, benchmark, hasFuzz, coverageRedesign bool) error {
	testDescr, err := parseTestSources(sources)
	if err != nil {
		return err
	}
	testDescr.Coverage = coverage
	testDescr.CoverVars = coverVars
	if len(testDescr.TestFunctions) > 0 || len(testDescr.BenchFunctions) > 0 || len(testDescr.Examples) > 0 || len(testDescr.FuzzFunctions) > 0 || testDescr.Main != "" {
		// Can't set this if there are no test functions, it'll be an unused import.
		if coverageRedesign {
			testDescr.Imports = []string{fmt.Sprintf("%s \"%s\"", testDescr.Package, testPackage)}
		} else {
			testDescr.Imports = extraImportPaths(testPackage, testDescr.Package, testDescr.CoverVars)
		}
	}

	testDescr.Benchmark = benchmark
	testDescr.HasFuzz = hasFuzz

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	// This might be consumed by other things.
	fmt.Printf("Package: %s\n", testDescr.Package)

	if coverageRedesign {
		return testMainTmpl.Execute(f, testDescr)
	}
	return oldTestMainTmpl.Execute(f, testDescr)
}

func extraImportPaths(testPackage, alias string, coverVars []CoverVar) []string {
	ret := make([]string, 0, len(coverVars)+1)
	ret = append(ret, fmt.Sprintf("%s \"%s\"", alias, testPackage))

	for i, v := range coverVars {
		name := fmt.Sprintf("_cover%d", i)
		coverVars[i].ImportName = name
		ret = append(ret, fmt.Sprintf("%s \"%s\"", name, v.ImportPath))
	}
	return ret
}

// parseTestSources parses the test sources and returns the package and set of test functions in them.
func parseTestSources(sources []string) (testDescr, error) {
	descr := testDescr{}
	for _, source := range sources {
		f, err := parser.ParseFile(token.NewFileSet(), source, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %s\n", source, err)
			return descr, err
		}
		descr.Package = f.Name.Name
		// If we're testing main, we will get errors from it clashing with func main.
		if descr.Package == "main" {
			descr.Package = "_main"
		}
		for _, d := range f.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv == nil {
				name := fd.Name.String()
				if isTestMain(fd) {
					descr.Main = name
				} else if isTest(fd, 1, name, "Test") {
					descr.TestFunctions = append(descr.TestFunctions, name)
				} else if isTest(fd, 1, name, "Benchmark") {
					descr.BenchFunctions = append(descr.BenchFunctions, name)
				} else if isTest(fd, 1, name, "Fuzz") {
					descr.FuzzFunctions = append(descr.FuzzFunctions, name)
				}
			}
		}
		// Get doc to find the examples for us :)
		descr.Examples = append(descr.Examples, doc.Examples(f)...)
	}
	return descr, nil
}

// isTestMain returns true if fn is a TestMain(m *testing.M) function.
// Copied from Go sources.
func isTestMain(fn *ast.FuncDecl) bool {
	if fn.Name.String() != "TestMain" ||
		fn.Type.Results != nil && len(fn.Type.Results.List) > 0 ||
		fn.Type.Params == nil ||
		len(fn.Type.Params.List) != 1 ||
		len(fn.Type.Params.List[0].Names) > 1 {
		return false
	}
	ptr, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	// We can't easily check that the type is *testing.M
	// because we don't know how testing has been imported,
	// but at least check that it's *M or *something.M.
	if name, ok := ptr.X.(*ast.Ident); ok && name.Name == "M" {
		return true
	}
	if sel, ok := ptr.X.(*ast.SelectorExpr); ok && sel.Sel.Name == "M" {
		return true
	}
	return false
}

// isTest returns true if the given function looks like a test.
// Copied from Go sources.
func isTest(fd *ast.FuncDecl, argLen int, name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) || fd.Recv != nil || len(fd.Type.Params.List) != argLen {
		return false
	} else if len(name) == len(prefix) { // "Test" is ok
		return true
	}

	rune, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(rune)
}

// testMainTmpl is the template for our test main, copied from Go's builtin one.
// Some bits are excluded because we don't support them and/or do them differently.
var testMainTmpl = template.Must(template.New("main").Parse(`
package main

import (
	_gostdlib_io "io"
	_gostdlib_os "os"
	{{if not .Benchmark}}_gostdlib_strings "strings"{{end}}
	_gostdlib_testing "testing"
	_gostdlib_testdeps "testing/internal/testdeps"

{{if .Coverage}}
	_ "runtime/coverage"
	_ "unsafe"
{{end}}

{{range .Imports}}
	{{.}}
{{end}}
)

var tests = []_gostdlib_testing.InternalTest{
{{range .TestFunctions}}
	{"{{.}}", {{$.Package}}.{{.}}},
{{end}}
}
var examples = []_gostdlib_testing.InternalExample{
{{range .Examples}}
	{"{{.Name}}", {{$.Package}}.Example{{.Name}}, {{.Output | printf "%q"}}, {{.Unordered}}},
{{end}}
}

var benchmarks = []_gostdlib_testing.InternalBenchmark{
{{range .BenchFunctions}}
	{"{{.}}", {{$.Package}}.{{.}}},
{{end}}
}

var fuzzTargets = []_gostdlib_testing.InternalFuzzTarget{
{{ range .FuzzFunctions }}
	{"{{.}}", {{$.Package}}.{{.}}},
{{ end }}
}

{{if .Coverage}}
//go:linkname runtime_coverage_processCoverTestDir runtime/coverage.processCoverTestDir
func runtime_coverage_processCoverTestDir(dir string, cfile string, cmode string, cpkgs string) error

//go:linkname testing_registerCover2 testing.registerCover2
func testing_registerCover2(mode string, tearDown func(coverprofile string, gocoverdir string) (string, error))

//go:linkname runtime_coverage_markProfileEmitted runtime/coverage.markProfileEmitted
func runtime_coverage_markProfileEmitted(val bool)

func coverTearDown(coverprofile string, gocoverdir string) (string, error) {
	var err error
	if gocoverdir == "" {
		gocoverdir, err = _gostdlib_os.MkdirTemp("", "gocoverdir")
		if err != nil {
			return "error setting GOCOVERDIR: bad os.MkdirTemp return", err
		}
		defer _gostdlib_os.RemoveAll(gocoverdir)
	}
	runtime_coverage_markProfileEmitted(true)
	if err := runtime_coverage_processCoverTestDir(gocoverdir, coverprofile, "set", ""); err != nil {
		return "error generating coverage report", err
	}
	return "", nil
}
{{end}}

var testDeps = _gostdlib_testdeps.TestDeps{}

func internalMain() int {
    if resultsFile := _gostdlib_os.Getenv("RESULTS_FILE"); resultsFile != "" {
        f, err := _gostdlib_os.Create(resultsFile)
        if err != nil {
            panic(err)
        }
        old := _gostdlib_os.Stdout
        mw := _gostdlib_io.MultiWriter(f, old)
        r, w, err := _gostdlib_os.Pipe()
        if err != nil {
            panic(err)
        }
        _gostdlib_os.Stdout = w
        done := make(chan struct{})
        go func() {
           _gostdlib_io.Copy(mw, r)
           done <- struct{}{}
         }()
        defer func() {
            w.Close()
            <-done
            r.Close()
            f.Close()
            _gostdlib_os.Stdout = old
        }()
    }
{{if .Coverage}}
    coverfile := _gostdlib_os.Getenv("COVERAGE_FILE")
    args := []string{_gostdlib_os.Args[0], "-test.v", "-test.coverprofile", coverfile}
	testing_registerCover2("set", coverTearDown)
{{else}}
    args := []string{_gostdlib_os.Args[0], "-test.v"}
{{end}}
{{if not .Benchmark}}
    testVar := _gostdlib_os.Getenv("TESTS")
    if testVar != "" {
		testVar = _gostdlib_strings.ReplaceAll(testVar, " ", "|")
		args = append(args, "-test.run", testVar)
    }
    _gostdlib_os.Args = append(args, _gostdlib_os.Args[1:]...)
	m := _gostdlib_testing.MainStart(testDeps, tests, nil, fuzzTargets, examples)
{{else}}
	args = append(args, "-test.bench", ".*")
	_gostdlib_os.Args = append(args, _gostdlib_os.Args[1:]...)
	m := _gostdlib_testing.MainStart(testDeps, nil, benchmarks, fuzzTargets, nil)
{{end}}
{{if .Main}}
	{{.Package}}.{{.Main}}(m)
    return 0
{{else}}
	return m.Run()
{{end}}
}

func main() {
	_gostdlib_os.Exit(internalMain())
}
`))

var oldTestMainTmpl = template.Must(template.New("oldmain").Parse(`
package main

import (
	_gostdlib_os "os"
	{{if not .Benchmark}}_gostdlib_strings "strings"{{end}}
	_gostdlib_testing "testing"
	_gostdlib_testdeps "testing/internal/testdeps"

{{range .Imports}}
	{{.}}
{{end}}
)

var tests = []_gostdlib_testing.InternalTest{
{{range .TestFunctions}}
	{"{{.}}", {{$.Package}}.{{.}}},
{{end}}
}
var examples = []_gostdlib_testing.InternalExample{
{{range .Examples}}
	{"{{.Name}}", {{$.Package}}.Example{{.Name}}, {{.Output | printf "%q"}}, {{.Unordered}}},
{{end}}
}

var benchmarks = []_gostdlib_testing.InternalBenchmark{
{{range .BenchFunctions}}
	{"{{.}}", {{$.Package}}.{{.}}},
{{end}}
}

{{ if .HasFuzz }}
var fuzzTargets = []_gostdlib_testing.InternalFuzzTarget{
{{ range .FuzzFunctions }}
	{"{{.}}", {{$.Package}}.{{.}}},
{{ end }}
}
{{ end }}

{{if .Coverage}}

// Only updated by init functions, so no need for atomicity.
var (
	coverCounters = make(map[string][]uint32)
	coverBlocks = make(map[string][]_gostdlib_testing.CoverBlock)
)

func init() {
	{{range $i, $c := .CoverVars}}
		{{if $c.ImportName }}
			coverRegisterFile({{printf "%q" $c.File}}, {{$c.ImportName}}.{{$c.Var}}.Count[:], {{$c.ImportName}}.{{$c.Var}}.Pos[:], {{$c.ImportName}}.{{$c.Var}}.NumStmt[:])
		{{end}}
	{{end}}
}

func coverRegisterFile(fileName string, counter []uint32, pos []uint32, numStmts []uint16) {
	if 3*len(counter) != len(pos) || len(counter) != len(numStmts) {
		panic("coverage: mismatched sizes")
	}
	if coverCounters[fileName] != nil {
		// Already registered.
		return
	}
	coverCounters[fileName] = counter
	block := make([]_gostdlib_testing.CoverBlock, len(counter))
	for i := range counter {
		block[i] = _gostdlib_testing.CoverBlock{
			Line0: pos[3*i+0],
			Col0: uint16(pos[3*i+2]),
			Line1: pos[3*i+1],
			Col1: uint16(pos[3*i+2]>>16),
			Stmts: numStmts[i],
		}
	}
	coverBlocks[fileName] = block
}
{{end}}

var testDeps = _gostdlib_testdeps.TestDeps{}

func main() {
{{if .Coverage}}
	_gostdlib_testing.RegisterCover(_gostdlib_testing.Cover{
		Mode: "set",
		Counters: coverCounters,
		Blocks: coverBlocks,
		CoveredPackages: "",
	})
    coverfile := _gostdlib_os.Getenv("COVERAGE_FILE")
    args := []string{_gostdlib_os.Args[0], "-test.v", "-test.coverprofile", coverfile}
{{else}}
    args := []string{_gostdlib_os.Args[0], "-test.v"}
{{end}}
{{if not .Benchmark}}
    testVar := _gostdlib_os.Getenv("TESTS")
    if testVar != "" {
		testVar = _gostdlib_strings.ReplaceAll(testVar, " ", "|")
		args = append(args, "-test.run", testVar)
    }
    _gostdlib_os.Args = append(args, _gostdlib_os.Args[1:]...)
	m := _gostdlib_testing.MainStart(testDeps, tests, nil,{{ if .HasFuzz }} fuzzTargets,{{ end }} examples)
{{else}}
	args = append(args, "-test.bench", ".*")
	_gostdlib_os.Args = append(args, _gostdlib_os.Args[1:]...)
	m := _gostdlib_testing.MainStart(testDeps, nil, benchmarks,{{ if .HasFuzz }} fuzzTargets,{{ end }} nil)
{{end}}

{{if .Main}}
	{{.Package}}.{{.Main}}(m)
{{else}}
	_gostdlib_os.Exit(m.Run())
{{end}}
}
`))

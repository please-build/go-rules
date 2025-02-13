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
	Imports        []string
	Benchmark      bool
}

// WriteTestMain templates a test main file from the given sources to the given output file.
func WriteTestMain(testPackage string, sources []string, output string, benchmark bool) error {
	testDescr, err := parseTestSources(sources)
	if err != nil {
		return err
	}
	if len(testDescr.TestFunctions) > 0 || len(testDescr.BenchFunctions) > 0 || len(testDescr.Examples) > 0 || len(testDescr.FuzzFunctions) > 0 || testDescr.Main != "" {
		testDescr.Imports = []string{fmt.Sprintf("%s \"%s\"", testDescr.Package, testPackage)}
	}

	testDescr.Benchmark = benchmark

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	// This might be consumed by other things.
	fmt.Printf("Package: %s\n", testDescr.Package)

	return testMainTmpl.Execute(f, testDescr)
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
	_gostdlib_os "os"
	{{ if not .Benchmark }}_gostdlib_strings "strings"{{ end }}
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

var fuzzTargets = []_gostdlib_testing.InternalFuzzTarget{
{{ range .FuzzFunctions }}
	{"{{.}}", {{$.Package}}.{{.}}},
{{ end }}
}

var testDeps = _gostdlib_testdeps.TestDeps{}

func internalMain() int {
	args := []string{_gostdlib_os.Args[0], "-test.v"}

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

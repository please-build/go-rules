package test

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTestSources(t *testing.T) {
	descr, err := parseTestSources([]string{"tools/please_go/test/test_data/test/example_test.go"})
	assert.NoError(t, err)
	assert.Equal(t, "buildgo", descr.Package)
	assert.Equal(t, "", descr.Main)
	functions := []string{
		"TestReadPkgdef",
		"TestReadCopiedPkgdef",
	}
	assert.Equal(t, functions, descr.TestFunctions)
}

func TestParseTestSourcesWithMain(t *testing.T) {
	descr, err := parseTestSources([]string{"tools/please_go/test/test_data/main/example_test_main.go"})
	assert.NoError(t, err)
	assert.Equal(t, "parse", descr.Package)
	assert.Equal(t, "TestMain", descr.Main)
	functions := []string{
		"TestParseSourceBuildLabel",
		"TestParseSourceRelativeBuildLabel",
		"TestParseSourceFromSubdirectory",
		"TestParseSourceFromOwnedSubdirectory",
		"TestParseSourceWithParentPath",
		"TestParseSourceWithAbsolutePath",
		"TestAddTarget",
	}
	assert.Equal(t, functions, descr.TestFunctions)
}

func TestParseTestSourcesSkipsExamplesWithoutOutput(t *testing.T) {
	descr, err := parseTestSources([]string{"tools/please_go/test/test_data/example/example_example_test.go"})
	assert.NoError(t, err)
	names := make([]string, 0, len(descr.Examples))
	for _, e := range descr.Examples {
		names = append(names, e.Name)
	}
	// Same rule as `go test`: an example is only run when it has an output comment, even an empty one.
	assert.Equal(t, []string{"EmptyOutput", "WithOutput"}, names)
}

func TestParseTestSourcesFailsGracefully(t *testing.T) {
	_, err := parseTestSources([]string{"wibble"})
	assert.Error(t, err)
}

func TestWriteTestMain(t *testing.T) {
	err := WriteTestMain("test_pkg", []string{"tools/please_go/test/test_data/test/example_test.go"}, "test.go", false)
	assert.NoError(t, err)
	// It's not really practical to assert the contents of the file in great detail.
	// We'll do the obvious thing of asserting that it is valid Go source.
	f, err := parser.ParseFile(token.NewFileSet(), "test.go", nil, 0)
	assert.NoError(t, err)
	assert.Equal(t, "main", f.Name.Name)
}

func TestWriteTestMainSkipsExamplesWithoutOutput(t *testing.T) {
	err := WriteTestMain("test_pkg", []string{"tools/please_go/test/test_data/example/example_example_test.go"}, "test.go", false)
	assert.NoError(t, err)
	f, err := parser.ParseFile(token.NewFileSet(), "test.go", nil, 0)
	assert.NoError(t, err)
	assert.Equal(t, "main", f.Name.Name)

	test, err := ioutil.ReadFile("test.go")
	assert.NoError(t, err)
	assert.Contains(t, string(test), "ExampleWithOutput")
	assert.Contains(t, string(test), "ExampleEmptyOutput")
	assert.NotContains(t, string(test), "ExampleNoOutput")
}

func TestWriteTestMainWithBenchmark(t *testing.T) {
	err := WriteTestMain("test_package", []string{"tools/please_go/test/test_data/bench/example_benchmark_test.go"}, "test.go", true)
	assert.NoError(t, err)
	// It's not really practical to assert the contents of the file in great detail.
	// We'll do the obvious thing of asserting that it is valid Go source.
	f, err := parser.ParseFile(token.NewFileSet(), "test.go", nil, 0)
	assert.NoError(t, err)
	assert.Equal(t, "main", f.Name.Name)

	test, err := ioutil.ReadFile("test.go")
	assert.NoError(t, err)
	assert.Contains(t, string(test), "BenchmarkExample")
}

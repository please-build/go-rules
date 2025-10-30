package benchmark

import (
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBenchDuration(t *testing.T) {
	r := require.New(t)

	dataCmd, _ := os.LookupEnv("DATA")
	cmd := exec.Command(dataCmd)

	start := time.Now()
	out, err := cmd.Output()
	r.NoError(err)

	r.Greater(int64(time.Since(start)), int64(time.Second), "Benchmark run was too quick to have actually run the tests")

	r.Contains(string(out), "BenchmarkOneSecWait")
	r.Contains(string(out), "Benchmark100msWait")
}

func TestExternalIsAddedToLabelsForExternalBenchmark(t *testing.T) {
	expectedTarget := "//test/benchmark:benchmark"

	output, err := exec.Command("plz", "query", "alltargets", "-i", "external").Output()
	require.NoError(t, err)

	outputLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	assert.True(t,
		slices.Contains(outputLines, expectedTarget),
		"Did not find %s in targets with the label 'external' even though external is set to True. Matching targets: %s",
		expectedTarget,
		outputLines,
	)
}

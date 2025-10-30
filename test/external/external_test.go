package external_test

import (
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/please-build/go-rules/test/external"
)

func TestAnswer(t *testing.T) {
	assert.Equal(t, 42, GetAnswer())
	assert.Equal(t, "var", GetVar())
	assert.Equal(t, "var1 var2", GetVar2())
}

func TestGitShow(t *testing.T) {
	const featureAdded = 1544432014
	lastCommitTime, err := strconv.ParseInt(GetExecGitShow(), 10, 64)
	if !assert.NoError(t, err) {
		assert.Fail(t, "unable to parse time")
	}
	assert.True(t, lastCommitTime > featureAdded, "git_show(): time went backwards")
}

func TestExecGitState(t *testing.T) {
	assert.Contains(t, GetExecGitState(), "shiny", "git_state(): failed")
}

func TestExecGitCommit(t *testing.T) {
	t.Skip("Failing on Alpine currently")
	assert.Len(t, GetExecGitCommit(), 40, "git_commit() length wrong")
}

func TestExecGitBranch(t *testing.T) {
	t.Skip("Failing on Alpine currently")
	assert.True(t, len(GetExecGitBranchFull()) > len(GetExecGitBranchShort()), "git_branch() lengths inconsistent")
	assert.Regexp(t, "^refs/", GetExecGitBranchFull())
	assert.Contains(t, GetExecGitBranchFull(), GetExecGitBranchShort(), "short branch should be in full branch")
}

func TestExternalIsAddedToLabelsForExternalTest(t *testing.T) {
	expectedTarget := "//test/external:external_test"

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

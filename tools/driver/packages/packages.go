package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/types"
	"os"
	"os/exec"

	"github.com/peterebden/go-cli-init/v5/logging"
	"golang.org/x/tools/go/packages"
)

var log = logging.MustGetLogger()

// DriverRequest is copied from go/packages; it's the format of config that we get on stdin.
type DriverRequest struct {
	Mode       packages.LoadMode `json:"mode"`
	Env        []string          `json:"env"`
	BuildFlags []string          `json:"build_flags"`
	Tests      bool              `json:"tests"`
	Overlay    map[string][]byte `json:"overlay"`
}

// DriverResponse is copied from go/packages; it's our response about packages we've loaded.
type DriverResponse struct {
	NotHandled bool
	Sizes      *types.StdSizes
	Roots      []string `json:",omitempty"`
	Packages   []*packages.Package
	GoVersion  int
}

// Load reads a set of packages and returns information about them.
// Most of the request structure isn't honoured at the moment.
func Load(req *DriverRequest, files []string) (*DriverResponse, error) {
	// We run Please as a subprocess to get the answers here.
	// A cooler way of handling this in future would be to do this in-process; for that we'd
	// need to define the SDK we keep talking about as a supported programmatic interface.
	whatinputs := exec.Command("plz", append([]string{"query", "whatinputs"}, files...)...)
	whatinputs.Stderr = &bytes.Buffer{}
	targetList, err := whatinputs.Output()
	if err != nil {
		return nil, handleSubprocessErr(whatinputs, err)
	}
	// We can't pipe whatinputs directly into plz build because we need the target list to use again.
	build := exec.Command("plz", "build", "-")
	build.Stdin = bytes.NewReader(targetList)
	build.Stderr = &bytes.Buffer{}
	if err := build.Run(); err != nil {
		return nil, handleSubprocessErr(build, err)
	}

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	defer w.Close()

	deps := exec.Command("plz", "query", "deps", "-")
	deps.Stdin = bytes.NewReader(targetList)
	deps.Stderr = &bytes.Buffer{}
	deps.Stdout = w
	if err := deps.Start(); err != nil {
		return nil, err
	}
	print := exec.Command("plz", "query", "print", "--json", "-")
	print.Stdin = r
	print.Stderr = &bytes.Buffer{}
	print.Stdout = &bytes.Buffer{}
	if err := print.Start(); err != nil {
		return nil, err
	}
	if err := deps.Wait(); err != nil {
		return nil, handleSubprocessErr(deps, err)
	}
	if err := print.Wait(); err != nil {
		return nil, handleSubprocessErr(print, err)
	}
	targets := map[string]*buildTarget{}
	if err := json.Unmarshal(print.Stdout.(*bytes.Buffer).Bytes(), targets); err != nil {
		return nil, err
	}
	return &DriverResponse{}, fmt.Errorf("Not handled: %s ", targets)
}

// buildTarget is a minimal version of the target output structure from `plz query print`
type buildTarget struct {
	Deps   []string    `json:"deps"`
	Labels []string    `json:"labels"`
	Srcs   interface{} `json:"srcs"`
}

func handleSubprocessErr(cmd *exec.Cmd, err error) error {
	return fmt.Errorf("%s Stdout:\n%s", err, cmd.Stderr.(*bytes.Buffer).String())
}

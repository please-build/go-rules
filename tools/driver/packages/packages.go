package packages

import (
	"fmt"
	"go/types"

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
func Load(req *DriverRequest, packages []string) (*DriverResponse, error) {
	return &DriverResponse{}, fmt.Errorf("Not handled")
}

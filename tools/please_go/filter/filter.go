package filter

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)


func Filter(tags []string, srcs []string) {
	ctxt := build.Default
	ctxt.BuildTags = tags

	for _, f := range srcs {
		dir, file := filepath.Split(f)

		// MatchFile skips _ prefixed files by default, assuming they're editor
		// temporary files - but we need to include cgo generated files.
		if strings.HasPrefix(file, "_cgo_") {
			fmt.Println(f)
			continue
		}

		ok, err := ctxt.MatchFile(dir, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking %v: %v\n", f, err)
			os.Exit(1)
		}

		if ok {
			fmt.Println(f)
		}
	}
}
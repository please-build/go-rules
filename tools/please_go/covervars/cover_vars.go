package covervars

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
)

func GenCoverVars(w io.Writer, importPath string, srcs []string) {
	for _, src := range srcs {
		if _, err := w.Write([]byte(coverVar(src, importPath))); err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func coverVar(src, importPath string) string {
	baseName := filepath.Base(src)
	return fmt.Sprintf("%s=GoCover_%s_go\n", importPath, strings.TrimSuffix(baseName, ".go"))
}

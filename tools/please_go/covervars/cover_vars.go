package covervars

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
)

// GenCoverVars writes the coverage variable information to w
func GenCoverVars(w io.Writer, importPath string, srcs []string) {
	for _, src := range srcs {
		if _, err := w.Write([]byte(coverVar(src, importPath))); err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func coverVar(src, importPath string) string {
	baseName := filepath.Base(src)
	return fmt.Sprintf("%s=GoCover_%s\n", importPath, strings.ReplaceAll(baseName, ".", "_"))
}

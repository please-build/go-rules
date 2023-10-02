package embed

import "embed"

//go:embed hello.txt
var hello string

//go:embed subdir
var subdir embed.FS

//go:embed all:subdir
var subdirAll embed.FS

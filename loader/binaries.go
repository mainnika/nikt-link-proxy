package loader

import (
	"embed"
)

// Static go-templates source
//go:embed binary/loader.js
var binariesFS embed.FS

const (
	LOADER_BIN = "binary/loader.js"
)

func init() {
	register(&binariesFS, LOADER_BIN)
}

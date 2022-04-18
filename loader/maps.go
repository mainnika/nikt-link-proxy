//go:build maps
// +build maps

package loader

import "embed"

// Static go-templates source
//go:embed binary/loader.js.map
var mapsFS embed.FS

const (
	LOADER_BIN_MAP = "binary/loader.js.map"
)

func init() {
	register(&mapsFS, LOADER_BIN_MAP)
	return
}

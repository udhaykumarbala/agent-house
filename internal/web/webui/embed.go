// Package webui embeds the Next.js static export (frontend/out) so the
// single Go binary still serves the whole UI — preserving the prebuilt-binary
// deploy model. Regenerate dist/ with scripts/build-ui.sh before `go build`.
package webui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

// FS returns the embedded export rooted at the dist directory.
func FS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err) // dist is embedded at compile time; this cannot fail
	}
	return sub
}

package ressam

import (
	"embed"
	"io/fs"
	"os"
)

const ClientPath = "web/dist"

//go:embed web/dist
var webFS embed.FS

// GetEmbeddedClientFS returns binary embedded filesystem.
func GetEmbeddedClientFS() fs.FS {
	f, err := fs.Sub(webFS, ClientPath)
	if err != nil {
		panic(err)
	}
	return f
}

// GetClientFS returns os implementation of fs.FS. Used only for local testing.
func GetClientFS() fs.FS {
	return os.DirFS(ClientPath)
}

package ressam

import (
	"embed"
	"io/fs"
)

//go:embed web/dist
var webFS embed.FS

func GetClientFS() fs.FS {
	f, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		panic(err)
	}
	return f
}

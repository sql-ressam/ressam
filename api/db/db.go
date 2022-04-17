package db

import (
	"context"
	"fmt"

	"github.com/sql-ressam/ressam/pg"
)

func init() {
	fmt.Println("test github actions")
}

// PgExporter can get information about PostgreSQL schemas.
type PgExporter interface {
	FetchDBInfo(ctx context.Context) (res pg.Info, err error)
}

// API handles database information requests.
type API struct {
	exporter PgExporter
}

// NewAPI returns new API instance.
func NewAPI(exporter PgExporter) API {
	return API{
		exporter,
	}
}

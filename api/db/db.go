package db

import (
	"context"

	"github.com/sql-ressam/ressam/database/pg"
)

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

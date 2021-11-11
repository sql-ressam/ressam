package db

import (
	"context"

	"github.com/sql-ressam/ressam/pg"
)

type PgExporter interface {
	GetDBInfo(ctx context.Context) (res pg.DBInfo, err error)
}

type API struct {
	exporter PgExporter
}

func NewAPI(exporter PgExporter) API {
	return API{
		exporter,
	}
}

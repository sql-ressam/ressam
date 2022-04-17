package pg

import (
	"github.com/sql-ressam/ressam/db"
)

// DBInfo is complete information about DBMS.
type DBInfo struct {
	Schemes []Scheme `json:"schemes"`
}

// Scheme is a PostgreSQL scheme.
type Scheme struct {
	Name          string            `json:"name"`
	Tables        []db.Table        `json:"tables"`
	Relationships []db.Relationship `json:"relationships"`
}

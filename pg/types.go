package pg

import (
	"github.com/sql-ressam/ressam/db"
)

// Info is complete information about DBMS.
type Info struct {
	Schemes []Scheme `json:"schemes"`
}

// Scheme is a PostgreSQL scheme.
type Scheme struct {
	Name          string            `json:"name"`
	Tables        []db.Table        `json:"tables"`
	Relationships []db.Relationship `json:"relationships"`
}

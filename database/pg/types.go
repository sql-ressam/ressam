package pg

import (
	"github.com/sql-ressam/ressam/database"
)

// Info is complete information about DBMS.
type Info struct {
	Schemes []Scheme `json:"schemes"`
}

// Scheme is a PostgreSQL scheme.
type Scheme struct {
	Name          string                  `json:"name"`
	Tables        []database.Table        `json:"tables"`
	Relationships []database.Relationship `json:"relationships"`
}

package pg

type DBInfo struct {
	Schemes []Scheme `json:"schemes"`
}

// Column contains info about a PostgreSQL column.
type Column struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Nullable       bool   `json:"nullable"`
	ColumnPosition int32  `json:"-"`

	DefaultValue *string `json:"omitempty,defaultValue"`
	Precision    *int    `json:"precision"`
}

// Table represents the PostgreSQL table.
type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

// ColumnInfo represents minimal info about a column.
type ColumnInfo struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

// Relationship represents Relationship between tables.
type Relationship struct {
	Name string     `json:"name"`
	From ColumnInfo `json:"from"`
	To   ColumnInfo `json:"to"`
}

// Scheme is PostgreSQL scheme.
type Scheme struct {
	Name          string         `json:"name"`
	Tables        []Table        `json:"tables"`
	Relationships []Relationship `json:"relationships"`
}

type ByColumnPosition []Column

func (c ByColumnPosition) Len() int {
	return len(c)
}

func (c ByColumnPosition) Less(i, j int) bool {
	return c[i].ColumnPosition < c[j].ColumnPosition
}

func (c ByColumnPosition) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

package db

// Table is a database table.
type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

// Column is a database db.
type Column struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Nullable       bool   `json:"nullable"`
	ColumnPosition int32  `json:"-"`

	DefaultValue *string `json:"omitempty,defaultValue"`
	Precision    *int    `json:"precision"`
}

// ColumnInfo is the minimum information about the db.
type ColumnInfo struct {
	Table     string `json:"table"`
	Column    string `json:"column"`
	IsVirtual bool   `json:"isVirtual"`
}

// Relationship is a relationship between tables.
type Relationship struct {
	Name string     `json:"name"`
	From ColumnInfo `json:"from"`
	To   ColumnInfo `json:"to"`
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
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

// ByColumnPosition implement sort.Interface for sorting columns by position.
type ByColumnPosition []Column

// Len returns the slice len.
func (c ByColumnPosition) Len() int {
	return len(c)
}

// Less compares i, j items.
func (c ByColumnPosition) Less(i, j int) bool {
	return c[i].ColumnPosition < c[j].ColumnPosition
}

// Swap i, j items.
func (c ByColumnPosition) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

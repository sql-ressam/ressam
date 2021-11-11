package pg

type DBInfo struct {
	Schemes []Scheme `json:"schemes"`
}

type Column struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Nullable       bool   `json:"nullable"`
	ColumnPosition int32  `json:"-"`

	DefaultValue *string `json:"defaultValue"`
	Precision    *int    `json:"precision"`
}

type Key struct {
	Name string `json:"name"`
}

type Index struct {
	Name string `json:"name"`
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
	Keys    []Key    `json:"keys"`
	Indexes []Index  `json:"indexes"`
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

type Scheme struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

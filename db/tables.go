package db

// Tables slice of Table.
type Tables []Table

// VirtualForeignKeys returns all virtual foreign keys.
func (t Tables) VirtualForeignKeys() []ColumnInfo {
	return nil
}

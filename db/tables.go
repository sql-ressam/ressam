package db

type Tables []Table

func (t Tables) VirtualForeignKeys() []ColumnInfo {
	return nil
}

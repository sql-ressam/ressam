package pg

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
)

type Exporter struct {
	db *sql.DB
}

func NewExporter(db *sql.DB) Exporter {
	return Exporter{
		db: db,
	}
}

const pgYesValue = "YES"

func (e Exporter) GetDBInfo(ctx context.Context) (res DBInfo, err error) {
	rows, err := e.db.QueryContext(ctx, `SELECT table_schema,
		   table_name,
		   column_name,
		   ordinal_position,
		   column_default,
		   is_nullable,
		   udt_name,
		   coalesce(character_maximum_length, numeric_precision, datetime_precision, interval_precision) as precision
	FROM information_schema.columns
	where table_schema = 'public'`)
	if err != nil {
		return res, fmt.Errorf("get information schems: %w", err)
	}

	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}(rows)

	columnByTableByScheme := map[string]map[string][]Column{}

	for rows.Next() {
		var (
			scheme, table, nullable string
			col                     Column
		)

		if err := rows.Scan(&scheme, &table, &col.Name, &col.ColumnPosition, &col.DefaultValue, &nullable, &col.Type,
			&col.Precision); err != nil {
			return res, fmt.Errorf("scan: %w", err)
		}

		// convert YES_OR_NO type to boolean
		col.Nullable = nullable == pgYesValue

		if columnByTableByScheme[scheme] == nil {
			columnByTableByScheme[scheme] = map[string][]Column{}
		}

		columnByTableByScheme[scheme][table] = append(columnByTableByScheme[scheme][table], col)
	}

	for schemeName, tableNames := range columnByTableByScheme {
		var tables []Table

		for tableName, columns := range tableNames {
			sort.Sort(ByColumnPosition(columns))

			tables = append(tables, Table{
				Name:    tableName,
				Columns: columns,
			})
		}

		res.Schemes = append(res.Schemes, Scheme{
			Name:   schemeName,
			Tables: tables,
		})
	}

	return res, nil
}

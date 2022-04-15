package pg

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/lib/pq"
)

type Exporter struct {
	db *sql.DB
}

func NewExporter(db *sql.DB) *Exporter {
	return &Exporter{
		db: db,
	}
}

const (
	pgYesValue = "YES"
)

func (e *Exporter) GetDBInfo(ctx context.Context) (res DBInfo, err error) {
	res.Schemes, err = e.getSchemes(ctx)
	if err != nil {
		return res, fmt.Errorf("get schemes: %w", err)
	}

	schemes := make([]string, len(res.Schemes))
	for idx := 0; idx < len(res.Schemes); idx++ {
		schemes[idx] = res.Schemes[idx].Name
	}

	relationshipsByScheme, err := e.getTablesRelationships(ctx, schemes)
	if err != nil {
		return res, fmt.Errorf("get relationships: %w", err)
	}

	for schemeName, schemeRelationships := range relationshipsByScheme {
		for idx := range res.Schemes {
			if res.Schemes[idx].Name == schemeName {
				res.Schemes[idx].Relationships = schemeRelationships
			}
		}
	}

	return res, nil
}

func (e *Exporter) getSchemes(ctx context.Context) (res []Scheme, err error) {
	rows, err := e.db.QueryContext(ctx, `
		SELECT table_schema,
		   table_name,
		   column_name,
		   ordinal_position,
		   column_default,
		   is_nullable,
		   udt_name,
		   coalesce(character_maximum_length, numeric_precision, datetime_precision, interval_precision) AS precision
		FROM information_schema.columns
		WHERE table_schema = 'public'`)
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

		err = rows.Scan(
			&scheme, &table, &col.Name,
			&col.ColumnPosition, &col.DefaultValue,
			&nullable, &col.Type, &col.Precision,
		)
		if err != nil {
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

		res = append(res, Scheme{
			Name:   schemeName,
			Tables: tables,
		})
	}

	return res, err
}

func (e *Exporter) getTablesRelationships(ctx context.Context, schemes []string) (map[string][]Relationship, error) {
	rows, err := e.db.QueryContext(ctx, `
		SELECT tc.table_schema, tc.constraint_name, tc.table_name, kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name
		FROM
			information_schema.table_constraints AS tc
				JOIN information_schema.key_column_usage AS kcu
					 ON tc.constraint_name = kcu.constraint_name
				JOIN information_schema.constraint_column_usage AS ccu
					 ON ccu.constraint_name = tc.constraint_name
		WHERE constraint_type = 'FOREIGN KEY' AND tc.table_schema = ANY($1)
		ORDER BY 1, 2;`, pq.Array(schemes))
	if err != nil {
		return nil, fmt.Errorf("get tables relationships query: %w", err)
	}

	res := make(map[string][]Relationship, len(schemes))

	for rows.Next() {
		var (
			rel    Relationship
			scheme string
		)

		err = rows.Scan(&scheme, &rel.Name, &rel.From.Table, &rel.From.Column, &rel.To.Table, &rel.To.Column)
		if err != nil {
			return nil, fmt.Errorf("while scan tables relationships: %w", err)
		}

		res[scheme] = append(res[scheme], rel)
	}

	return res, nil
}

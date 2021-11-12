package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testConn *sql.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		log.Fatalln("TEST_DB_DSN environment variable is required")
	}

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln("open testConn", err.Error())
	}

	testConn = conn

	m.Run()
}

func TestExporter_GetDBInfo(t *testing.T) {
	exporter := NewExporter(testConn)

	t.Run("fetch database info without errors", func(t *testing.T) {
		info, err := exporter.GetDBInfo(context.Background())
		assert.Nil(t, err)

		infoMarshal, err := json.MarshalIndent(info, "", "\t")
		assert.Nil(t, err)
		t.Logf(string(infoMarshal))
	})

	t.Run("fetch info about test_default_values", func(t *testing.T) {
		info, err := exporter.GetDBInfo(context.Background())
		assert.Nil(t, err)

		expect := Table{
			Name: "test_default_values",
			Columns: []Column{
				{
					Name:           "id",
					Type:           "int8",
					Nullable:       false,
					ColumnPosition: 1,
					DefaultValue:   getStrRef("nextval('test_default_values_id_seq'::regclass)"),
					Precision:      getIntRef(64),
				},
				{
					Name:           "int_null",
					Type:           "int4",
					Nullable:       true,
					ColumnPosition: 2,
					DefaultValue:   nil,
					Precision:      getIntRef(32),
				},
				{
					Name:           "int_not_null",
					Type:           "int4",
					Nullable:       false,
					ColumnPosition: 3,
					DefaultValue:   nil,
					Precision:      getIntRef(32),
				},
				{
					Name:           "int_null_default_1",
					Type:           "int4",
					Nullable:       true,
					ColumnPosition: 4,
					DefaultValue:   getStrRef("1"),
					Precision:      getIntRef(32),
				},
				{
					Name:           "int_not_null_default_1",
					Type:           "int4",
					Nullable:       false,
					ColumnPosition: 5,
					DefaultValue:   getStrRef("1"),
					Precision:      getIntRef(32),
				},

				{
					Name:           "test_enum_null",
					Type:           "test_enum",
					Nullable:       true,
					ColumnPosition: 6,
					DefaultValue:   nil,
					Precision:      nil,
				},
				{
					Name:           "test_enum_not_null",
					Type:           "test_enum",
					Nullable:       false,
					ColumnPosition: 7,
					DefaultValue:   nil,
					Precision:      nil,
				},
				{
					Name:           "test_enum_null_default_first",
					Type:           "test_enum",
					Nullable:       true,
					ColumnPosition: 8,
					DefaultValue:   getStrRef("'first'::test_enum"),
					Precision:      nil,
				},
				{
					Name:           "test_enum_not_null_default_first",
					Type:           "test_enum",
					Nullable:       false,
					ColumnPosition: 9,
					DefaultValue:   getStrRef("'first'::test_enum"),
					Precision:      nil,
				},
			},
		}

		for _, table := range info.Schemes[0].Tables {
			if table.Name == "test_default_values" {
				assert.Equal(t, expect, table)
				break
			}
		}
	})
}

func getStrRef(s string) *string {
	return &s
}

func getIntRef(i int) *int {
	return &i
}

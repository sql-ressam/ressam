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

	"github.com/sql-ressam/ressam/internal/pkg/help"
)

var testConn *sql.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		log.Println("TEST_DB_DSN environment variable is required")
		os.Exit(2)
	}

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("open testConn", err.Error())
		os.Exit(2)
	}

	if err := conn.Ping(); err != nil {
		log.Println("ping:", err.Error())
		os.Exit(2)
	}

	testConn = conn

	os.Exit(m.Run())
}

func TestExporter_GetDBInfo(t *testing.T) {
	exporter := NewExporter(testConn)

	t.Run("fetch database info without errors", func(t *testing.T) {
		info, err := exporter.GetDBInfo(context.Background())
		assert.NoError(t, err)

		_, err = json.MarshalIndent(info, "", "\t")
		assert.NoError(t, err)
	})

	t.Run("fetch info about test_default_values", func(t *testing.T) {
		info, err := exporter.GetDBInfo(context.Background())
		assert.NoError(t, err)

		expect := Table{
			Name: "test_default_values",
			Columns: []Column{
				{
					Name:           "id",
					Type:           "int8",
					Nullable:       false,
					ColumnPosition: 1,
					DefaultValue:   help.Ref[string]("nextval('test_default_values_id_seq'::regclass)"),
					Precision:      help.Ref[int](64),
				},
				{
					Name:           "int_null",
					Type:           "int4",
					Nullable:       true,
					ColumnPosition: 2,
					DefaultValue:   nil,
					Precision:      help.Ref[int](32),
				},
				{
					Name:           "int_not_null",
					Type:           "int4",
					Nullable:       false,
					ColumnPosition: 3,
					DefaultValue:   nil,
					Precision:      help.Ref[int](32),
				},
				{
					Name:           "int_null_default_1",
					Type:           "int4",
					Nullable:       true,
					ColumnPosition: 4,
					DefaultValue:   help.Ref[string]("1"),
					Precision:      help.Ref[int](32),
				},
				{
					Name:           "int_not_null_default_1",
					Type:           "int4",
					Nullable:       false,
					ColumnPosition: 5,
					DefaultValue:   help.Ref[string]("1"),
					Precision:      help.Ref[int](32),
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
					DefaultValue:   help.Ref[string]("'first'::test_enum"),
					Precision:      nil,
				},
				{
					Name:           "test_enum_not_null_default_first",
					Type:           "test_enum",
					Nullable:       false,
					ColumnPosition: 9,
					DefaultValue:   help.Ref[string]("'first'::test_enum"),
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

	t.Run("check relationships", func(t *testing.T) {
		info, err := exporter.GetDBInfo(context.Background())
		assert.Nil(t, err)

		/*
			public,books_authors_author_id_fkey,books_authors,author_id,authors,id
			public,books_authors_book_id_fkey,books_authors,book_id,books,id
			public,books_users_book_id_fkey,books_users,book_id,books,id
			public,books_users_user_id_fkey,books_users,user_id,users,id
		*/
		expect := map[string][]Relationship{
			"public": {
				{
					Name: "books_authors_author_id_fkey",
					From: ColumnInfo{
						Table:  "books_authors",
						Column: "author_id",
					},
					To: ColumnInfo{
						Table:  "authors",
						Column: "id",
					},
				},
				{
					Name: "books_authors_book_id_fkey",
					From: ColumnInfo{
						Table:  "books_authors",
						Column: "book_id",
					},
					To: ColumnInfo{
						Table:  "books",
						Column: "id",
					},
				},
				{
					Name: "books_users_book_id_fkey",
					From: ColumnInfo{
						Table:  "books_users",
						Column: "book_id",
					},
					To: ColumnInfo{
						Table:  "books",
						Column: "id",
					},
				},
				{
					Name: "books_users_user_id_fkey",
					From: ColumnInfo{
						Table:  "books_users",
						Column: "user_id",
					},
					To: ColumnInfo{
						Table:  "users",
						Column: "id",
					},
				},
			},
		}

		for expectSchemeName, expectRelationships := range expect {
			for _, resultScheme := range info.Schemes {
				if resultScheme.Name == expectSchemeName {
					assert.Equal(t, expectRelationships, resultScheme.Relationships)
				}
			}
		}
	})
}

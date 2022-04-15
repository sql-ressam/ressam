package db

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sql-ressam/ressam/internal/pkg/help"
	"github.com/sql-ressam/ressam/pg"
)

func (a API) DBInfo(w http.ResponseWriter, r *http.Request) {
	info, err := a.exporter.GetDBInfo(r.Context())
	if err != nil {
		http.Error(w, "get db info: "+err.Error(), http.StatusInternalServerError)
	}

	res, err := json.Marshal(info)
	if err != nil {
		http.Error(w, "marshal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		log.Println("write:", err.Error())
	}
}

func FakeDBInfo(w http.ResponseWriter, r *http.Request) {
	fakeTable := pg.Table{
		Name: "test_default_values",
		Columns: []pg.Column{
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
			},
			{
				Name:           "test_enum_not_null",
				Type:           "test_enum",
				ColumnPosition: 7,
			},
			{
				Name:           "test_enum_null_default_first",
				Type:           "test_enum",
				Nullable:       true,
				ColumnPosition: 8,
				DefaultValue:   help.Ref[string]("'first'::test_enum"),
			},
			{
				Name:           "test_enum_not_null_default_first",
				Type:           "test_enum",
				ColumnPosition: 9,
				DefaultValue:   help.Ref[string]("'first'::test_enum"),
			},
		},
	}

	rels := []pg.Relationship{
		{
			Name: "books_authors_author_id_fkey",
			From: pg.ColumnInfo{
				Table:  "books_authors",
				Column: "author_id",
			},
			To: pg.ColumnInfo{
				Table:  "authors",
				Column: "id",
			},
		},
		{
			Name: "books_authors_book_id_fkey",
			From: pg.ColumnInfo{
				Table:  "books_authors",
				Column: "book_id",
			},
			To: pg.ColumnInfo{
				Table:  "books",
				Column: "id",
			},
		},
		{
			Name: "books_users_book_id_fkey",
			From: pg.ColumnInfo{
				Table:  "books_users",
				Column: "book_id",
			},
			To: pg.ColumnInfo{
				Table:  "books",
				Column: "id",
			},
		},
		{
			Name: "books_users_user_id_fkey",
			From: pg.ColumnInfo{
				Table:  "books_users",
				Column: "user_id",
			},
			To: pg.ColumnInfo{
				Table:  "users",
				Column: "id",
			},
		},
	}

	info := pg.DBInfo{
		Schemes: []pg.Scheme{
			{
				Name:          "public",
				Tables:        []pg.Table{fakeTable},
				Relationships: rels,
			},
		},
	}

	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Println(err.Error())
	}
}

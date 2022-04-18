package db

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sql-ressam/ressam/database"
	"github.com/sql-ressam/ressam/database/pg"
	"github.com/sql-ressam/ressam/pkg/help"
)

// DBInfo fetches info about database and writes the result.
func (a API) DBInfo(w http.ResponseWriter, r *http.Request) {
	info, err := a.exporter.FetchDBInfo(r.Context())
	if err != nil {
		http.Error(w, "get database info: "+err.Error(), http.StatusInternalServerError)
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

// FakeDBInfo returns PostgreSQL fake info. Used for local testing.
func FakeDBInfo(w http.ResponseWriter, _ *http.Request) {
	fakeTables := []database.Table{
		{
			Name: "books_users",
			Columns: []database.Column{
				{
					Name:      "user_id",
					Type:      "int8",
					Precision: help.Ref[int](64),
				},
				{
					Name:      "book_id",
					Type:      "int8",
					Precision: help.Ref[int](64),
				},
				{
					Name:         "receiving_date",
					Type:         "timestamp",
					DefaultValue: help.Ref[string]("now()"),
					Precision:    help.Ref[int](6),
				},
			},
		},
		{
			Name: "books_authors",
			Columns: []database.Column{
				{
					Name:      "author_id",
					Type:      "int8",
					Precision: help.Ref[int](64),
				},
				{
					Name:      "book_id",
					Type:      "int8",
					Precision: help.Ref[int](64),
				},
			},
		},
		{
			Name: "users",
			Columns: []database.Column{
				{
					Name:         "id",
					Type:         "int8",
					DefaultValue: help.Ref[string]("nextval('users_id_seq'::regclass)"),
					Precision:    help.Ref[int](64),
				},
				{
					Name:      "name",
					Type:      "varchar",
					Nullable:  true,
					Precision: help.Ref[int](255),
				},
			},
		},
		{
			Name: "books",
			Columns: []database.Column{
				{
					Name:         "id",
					Type:         "int8",
					DefaultValue: help.Ref[string]("nextval('books_id_seq'::regclass)"),
					Precision:    help.Ref[int](64),
				},
				{
					Name:      "name",
					Type:      "varchar",
					Nullable:  true,
					Precision: help.Ref[int](255),
				},
			},
		},
		{
			Name: "authors",
			Columns: []database.Column{
				{
					Name:         "id",
					Type:         "int8",
					DefaultValue: help.Ref[string]("nextval('authors_id_seq'::regclass)"),
					Precision:    help.Ref[int](64),
				},
				{
					Name:      "name",
					Type:      "varchar",
					Nullable:  true,
					Precision: help.Ref[int](255),
				},
			},
		},
	}

	rels := []database.Relationship{
		{
			Name: "books_authors_author_id_fkey",
			From: database.ColumnInfo{
				Table:  "books_authors",
				Column: "author_id",
			},
			To: database.ColumnInfo{
				Table:  "authors",
				Column: "id",
			},
		},
		{
			Name: "books_authors_book_id_fkey",
			From: database.ColumnInfo{
				Table:  "books_authors",
				Column: "book_id",
			},
			To: database.ColumnInfo{
				Table:  "books",
				Column: "id",
			},
		},
		{
			Name: "books_users_book_id_fkey",
			From: database.ColumnInfo{
				Table:  "books_users",
				Column: "book_id",
			},
			To: database.ColumnInfo{
				Table:  "books",
				Column: "id",
			},
		},
		{
			Name: "books_users_user_id_fkey",
			From: database.ColumnInfo{
				Table:  "books_users",
				Column: "user_id",
			},
			To: database.ColumnInfo{
				Table:  "users",
				Column: "id",
			},
		},
	}

	info := pg.Info{
		Schemes: []pg.Scheme{
			{
				Name:          "public",
				Tables:        fakeTables,
				Relationships: rels,
			},
		},
	}

	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Println(err.Error())
	}
}

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

func FakeDBInfo(w http.ResponseWriter, _ *http.Request) {
	fakeTables := []pg.Table{
		{
			Name: "books_users",
			Columns: []pg.Column{
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
			Columns: []pg.Column{
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
			Columns: []pg.Column{
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
			Columns: []pg.Column{
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
			Columns: []pg.Column{
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
				Tables:        fakeTables,
				Relationships: rels,
			},
		},
	}

	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Println(err.Error())
	}
}

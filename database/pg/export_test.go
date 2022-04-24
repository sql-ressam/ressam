//go:build integration

package pg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sql-ressam/ressam/database"
	"github.com/sql-ressam/ressam/pkg/help"
)

var expect = database.Table{
	Name: "test_default_values",
	Columns: []database.Column{
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

func TestExporter_GetDBInfo(t *testing.T) {
	type TestCase struct {
		Name  string
		PgTag string
	}

	testCases := []TestCase{
		{
			Name:  "postgres_v14",
			PgTag: "14",
		},
	}

	for _, tc := range testCases {
		db, release, err := preparePg(tc.PgTag)
		if err != nil {
			panic(err)
		}

		exporter := NewExporter(db)

		t.Run(tc.Name, func(t *testing.T) {
			ctx, stop := context.WithTimeout(context.Background(), time.Minute)
			info, err := exporter.FetchDBInfo(ctx)
			stop()
			require.NoError(t, err)
			assert.Len(t, info.Schemes, 1)

			for _, table := range info.Schemes[0].Tables {
				if table.Name == "test_default_values" {
					assert.Equal(t, expect, table)
					break
				}
			}
		})

		if err = release(); err != nil {
			panic(err)
		}
	}

	//exporter := NewExporter(testConn)
	//
	//t.Run("fetch database info without errors", func(t *testing.T) {
	//	info, err := exporter.FetchDBInfo(context.Background())
	//	assert.NoError(t, err)
	//
	//	_, err = json.MarshalIndent(info, "", "\t")
	//	assert.NoError(t, err)
	//})
	//

	//
	//t.Run("check relationships", func(t *testing.T) {
	//	info, err := exporter.FetchDBInfo(context.Background())
	//	assert.Nil(t, err)
	//
	//	/*
	//		public,books_authors_author_id_fkey,books_authors,author_id,authors,id
	//		public,books_authors_book_id_fkey,books_authors,book_id,books,id
	//		public,books_users_book_id_fkey,books_users,book_id,books,id
	//		public,books_users_user_id_fkey,books_users,user_id,users,id
	//	*/
	//	expect := map[string][]database.Relationship{
	//		"public": {
	//			{
	//				Name: "books_authors_author_id_fkey",
	//				From: database.ColumnInfo{
	//					Table:  "books_authors",
	//					Column: "author_id",
	//				},
	//				To: database.ColumnInfo{
	//					Table:  "authors",
	//					Column: "id",
	//				},
	//			},
	//			{
	//				Name: "books_authors_book_id_fkey",
	//				From: database.ColumnInfo{
	//					Table:  "books_authors",
	//					Column: "book_id",
	//				},
	//				To: database.ColumnInfo{
	//					Table:  "books",
	//					Column: "id",
	//				},
	//			},
	//			{
	//				Name: "books_users_book_id_fkey",
	//				From: database.ColumnInfo{
	//					Table:  "books_users",
	//					Column: "book_id",
	//				},
	//				To: database.ColumnInfo{
	//					Table:  "books",
	//					Column: "id",
	//				},
	//			},
	//			{
	//				Name: "books_users_user_id_fkey",
	//				From: database.ColumnInfo{
	//					Table:  "books_users",
	//					Column: "user_id",
	//				},
	//				To: database.ColumnInfo{
	//					Table:  "users",
	//					Column: "id",
	//				},
	//			},
	//		},
	//	}
	//
	//	for expectSchemeName, expectRelationships := range expect {
	//		for _, resultScheme := range info.Schemes {
	//			if resultScheme.Name == expectSchemeName {
	//				assert.Equal(t, expectRelationships, resultScheme.Relationships)
	//			}
	//		}
	//	}
	//})
}

func preparePg(tag string) (db *sql.DB, release func() error, err error) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not create a pool: %w", err)
	}

	const (
		repository = "postgres"
		user       = "postgres"
		password   = "postgres"
		host       = "localhost"
		driver     = "postgres"
		dbName     = "ressam"
	)
	env := []string{"POSTGRES_PASSWORD=" + password, "POSTGRES_USER=" + user, "POSTGRES_DB=" + dbName}
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: repository,
		Env:        env,
		Tag:        tag,
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}

	release = func() error {
		return pool.Purge(resource)
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", password, user, host,
		resource.GetPort("5432/tcp"), dbName)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(func() error {
		db, err = sql.Open(driver, dsn)
		if err != nil {
			log.Println("can't open sql connection:", err.Error())
			return err
		}
		return db.Ping()
	}); err != nil {
		if errRelease := release(); errRelease != nil {
			log.Println("can't remove container:", errRelease.Error())
		}
		return nil, nil, fmt.Errorf("can't open connection: %w", err)
	}

	return db, release, nil
}

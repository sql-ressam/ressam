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
}

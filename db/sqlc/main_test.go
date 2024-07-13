package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/alikhanMuslim/simpleBank/util"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

var testqueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot open config")
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("err to connect database", err)
	}

	testqueries = New(testDB)

	os.Exit(m.Run())
}

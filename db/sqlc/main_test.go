package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	drivername = "postgres"
	dataSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testDB *sql.DB

var testqueries *Queries

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(drivername, dataSource)

	if err != nil {
		log.Fatal("err to connect database", err)
	}

	testqueries = New(testDB)

	os.Exit(m.Run())
}

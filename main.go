package main

import (
	"database/sql"
	"log"

	"github.com/alikhanMuslim/simpleBank/api"
	db "github.com/alikhanMuslim/simpleBank/db/sqlc"
	"github.com/alikhanMuslim/simpleBank/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannnot load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("err to connect database", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

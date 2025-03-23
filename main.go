package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/varsilias/simplebank/api"
	db "github.com/varsilias/simplebank/db/sqlc"
	"github.com/varsilias/simplebank/utils"
)

func main() {
	var err error
	config, err := utils.LoadConfig(".env")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBUrl)

	if err != nil {
		log.Fatal("could not connect to database: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)

	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("could not start server: ", err)
	}

}

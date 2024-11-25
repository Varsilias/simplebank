package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/varsilias/simplebank/api"
	db "github.com/varsilias/simplebank/db/sqlc"
	"github.com/varsilias/simplebank/utils"
	"log"
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
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("could not start server: ", err)
	}

}

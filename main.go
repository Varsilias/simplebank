package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/varsilias/simplebank/api"
	db "github.com/varsilias/simplebank/db/sqlc"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://simplebank:SimpleBank1234@localhost:5432/simplebank?sslmode=disable"
	serverAddress = "0.0.0.0:8000"
)

func main() {

	var err error
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("could not connect to database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("could not start server: ", err)
	}

}

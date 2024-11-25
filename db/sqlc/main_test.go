package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/varsilias/simplebank/utils"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../test.env")
	if err != nil {
		log.Fatal("TEST: could not load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBUrl)

	if err != nil {
		log.Fatal("could not connect to database: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}

func createTestUser() User {
	args := CreateUserParams{
		PublicID:  utils.RandomString(26),
		Firstname: utils.RandomName(),
		Lastname:  utils.RandomName(),
		Email:     utils.RandomEmail(),
		Password:  utils.RandomString(16),
		Salt:      utils.RandomString(14),
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	if err != nil {
		log.Fatal("Error creating test user: ", err)
	}

	return user
}

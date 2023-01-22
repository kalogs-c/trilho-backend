package models_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kalogsc/ego/seed"
	"github.com/kalogsc/ego/server"
	"github.com/kalogsc/ego/utils"
)

var serverInstance = server.Server{}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error loading env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	err := serverInstance.InstanciateDB(os.Getenv("TEST_DB_NAME"), utils.DB_MODE_TEST)
	if err != nil {
		fmt.Println("Cannot connect to the database")
		log.Fatal("This is the error: ", err)
	} else {
		fmt.Println("We are connected to the database")
	}

	seed.Load(serverInstance.DB)
}

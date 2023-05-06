package server_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kalogsc/trilho/models"
	"github.com/kalogsc/trilho/server"
)

var serverInstance = server.Server{}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error loading env %v\n", err)
	}
	database()

	os.Exit(m.Run())
}

func database() {
	err := serverInstance.InstanciateDB()
	if err != nil {
		fmt.Println("Cannot connect to the database")
		log.Fatal("This is the error: ", err)
	} else {
		fmt.Println("We are connected to the database")
	}

	serverInstance.DB.Debug().AutoMigrate(&models.User{})
	serverInstance.DB.Debug().AutoMigrate(&models.Transaction{})
}

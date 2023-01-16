package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kalogsc/ego/seed"
	"github.com/kalogsc/ego/server"
)

var serverInstance server.Server = server.Server{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env's. Err: %v", err)
	}

	if len(os.Args) > 1 && (os.Args[1] == "--test" || os.Args[1] == "-t") {
		serverInstance.Initialize(os.Getenv("TEST_DB_NAME"))
		seed.Load(serverInstance.DB)
	} else {
		serverInstance.Initialize(os.Getenv("DB_NAME"))
	}

	serverInstance.Run(":8080")
}

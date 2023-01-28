package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kalogsc/ego/seed"
	"github.com/kalogsc/ego/server"
	"github.com/kalogsc/ego/utils"
)

var serverInstance server.Server = server.Server{}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--test" || os.Args[1] == "-t") {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error getting env's. Err: %v", err)
		}
		serverInstance.Initialize(os.Getenv("TEST_DB_NAME"), utils.DB_MODE_TEST)
		seed.Load(serverInstance.DB)
	} else {
		serverInstance.Initialize(os.Getenv("DB_NAME"), utils.DB_MODE_PROD)
		seed.Load(serverInstance.DB)
	}

	serverInstance.Run(":8080")
}

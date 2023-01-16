package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kalogsc/ego/server"
)

var serverInstance server.Server = server.Server{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env's. Err: %v", err)
	}

	serverInstance.Initialize(os.Getenv("DB_NAME"))

	serverInstance.Run(":8080")
}

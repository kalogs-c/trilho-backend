package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kalogsc/trilho/seed"
	"github.com/kalogsc/trilho/server"
)

var serverInstance server.Server = server.Server{}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--test" || os.Args[1] == "-t") {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error getting env's. Err: %v", err)
		}
		seed.Load(serverInstance.DB)
	} 
	serverInstance.Initialize()
	serverInstance.Run(":8080")
}

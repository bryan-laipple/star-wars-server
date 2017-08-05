package main

import (
	"os"

	"github.com/bryan-laipple/star-wars-server/server"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = defaultPort
	}
	server.Start(port)
}

package main // важно: не package cmd

import (
	"log"

	"my-notes-app/server" // имя должно совпасть с module в go.mod
)

func main() {
	srv := server.New()
	log.Println("starting server on port 8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

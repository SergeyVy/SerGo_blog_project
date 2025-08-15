package main

import (
	"log"
	"os"

	"my-notes-app/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := server.New()
	log.Println("starting server on port", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

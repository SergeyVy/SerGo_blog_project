package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"my-notes-app/server" // проверь, что совпадает с module в go.mod
)

func main() {
	// 1) Подключаемся к БД
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()
	log.Println("DB connected")

	// 2) Порт для Render/локально
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 3) Запускаем твой сервер (в server.New уже есть /ping)
	srv := server.New()
	log.Println("starting server on port", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

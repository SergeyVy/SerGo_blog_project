package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"my-notes-app/config"
	"my-notes-app/internal/handlers"
	"my-notes-app/pkg/logger"
	"my-notes-app/server"
	"my-notes-app/storage"

	_ "github.com/lib/pq"
)

func main() {
	// Логгер
	logger.Init()

	// Конфиг
	cfg := config.MustLoad()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = cfg.DatabaseURL
		log.Println("Using database_url from config.yaml:", dsn)
	} else {
		log.Println("Using DATABASE_URL from environment:", dsn)
	}

	// Подключение к БД
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	log.Println("Connected to database")

	// Миграции
	store := storage.New(db)
	if err := store.Migrate(context.Background()); err != nil {
		log.Fatal("DB migrate error:", err)
	}

	// Хендлер с зависимостями
	h := handlers.NewHandler(store, cfg.JWTSecret)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.New(h)
	log.Println("Starting server on port", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

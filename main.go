package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/internal/server"
)

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		if os.Getenv("VERCEL") != "" || !fileExists("portfolio.db") {
			dbPath = "file:memdb?mode=memory&cache=shared"
		} else {
			dbPath = "portfolio.db"
		}
	}

	webhookURL := os.Getenv("WEBHOOK_URL")

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := repository.SeedData(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	srv, err := server.NewServer(db)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	log.Printf("Listening on http://localhost:%s (DB: %s)", port, dbPath)
	if webhookURL != "" {
		log.Printf("Webhook: %s", webhookURL)
	}

	if err := http.ListenAndServe(":"+port, srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("http server error: %v", err)
	}
}

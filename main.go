package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		if os.Getenv("VERCEL") != "" {
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

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: srv,
	}

	log.Printf("Listening on :%s (DB: %s)", port, dbPath)
	if webhookURL != "" {
		log.Printf("Webhook: %s", webhookURL)
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
}

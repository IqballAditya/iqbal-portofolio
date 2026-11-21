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
		dbPath = "portfolio.db"
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

	go func() {
		log.Printf("==================================================")
		log.Printf("  IQBAL Portfolio Platform - Single Binary Server")
		log.Printf("==================================================")
		log.Printf("  Listening on    : http://localhost:%s", port)
		log.Printf("  Database        : %s (pure-Go SQLite / WAL mode)", dbPath)
		if webhookURL != "" {
			log.Printf("  Webhook URL     : %s", webhookURL)
		}
		log.Printf("  Embedded Assets : Static CSS/JS & Templates")
		log.Printf("  Routes:")
		log.Printf("    - GET  /                     (Home Showcase Feed)")
		log.Printf("    - GET  /gateway              (Persona Switcher)")
		log.Printf("    - POST /gateway/select       (Select Persona)")
		log.Printf("    - GET  /search               (Search Matrix)")
		log.Printf("    - GET  /api/search           (Live Search API)")
		log.Printf("    - GET  /projects/{id}/modal  (Case Study Modal)")
		log.Printf("    - GET  /projects/{id}/player (Interactive Demo Player)")
		log.Printf("    - GET  /shortlist            (Shortlist & Engagement Tiers)")
		log.Printf("    - POST /shortlist/{id}       (Toggle Shortlist Item)")
		log.Printf("    - POST /inquiries            (Hire Me Contact Form)")
		log.Printf("    - GET  /healthz              (System Health Check)")
		log.Printf("==================================================")

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server stopped cleanly.")
}

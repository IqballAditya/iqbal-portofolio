package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"

	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/internal/server"
)

var (
	srv     *server.Server
	once    sync.Once
	initErr error
)

func initApp() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "file:memdb?mode=memory&cache=shared"
	}

	db, err := database.InitDB(dbPath)
	if err != nil {
		initErr = fmt.Errorf("InitDB failed: %w", err)
		return
	}

	if err := database.Migrate(db); err != nil {
		initErr = fmt.Errorf("Migrate failed: %w", err)
		return
	}

	if err := repository.SeedData(db); err != nil {
		initErr = fmt.Errorf("SeedData failed: %w", err)
		return
	}

	s, err := server.NewServer(db)
	if err != nil {
		initErr = fmt.Errorf("NewServer failed: %w", err)
		return
	}
	srv = s
}

func Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Panic in Handler: %v\nStack:\n%s", rec, debug.Stack())
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 Internal Server Error: %v\n\n%s", rec, debug.Stack())
		}
	}()

	once.Do(initApp)
	if initErr != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Initialization error: %v", initErr)
		return
	}

	// Resolve target path
	targetPath := "/"
	if p := r.URL.Query().Get("__path"); p != "" {
		targetPath = "/" + strings.TrimPrefix(p, "/")
	} else if r.URL.Path != "" && r.URL.Path != "/api" && r.URL.Path != "/api/index" {
		targetPath = r.URL.Path
	} else if orig := r.Header.Get("x-matched-path"); orig != "" && !strings.HasPrefix(orig, "/api") {
		targetPath = orig
	}

	r.URL.Path = targetPath
	srv.ServeHTTP(w, r)
}

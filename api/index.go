package handler

import (
	"net/http"
	"os"
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
		dbPath = ":memory:"
	}

	db, err := database.InitDB(dbPath)
	if err != nil {
		initErr = err
		return
	}

	if err := database.Migrate(db); err != nil {
		initErr = err
		return
	}

	if err := repository.SeedData(db); err != nil {
		initErr = err
		return
	}

	s, err := server.NewServer(db)
	if err != nil {
		initErr = err
		return
	}
	srv = s
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	if initErr != nil {
		http.Error(w, "Initialization error: "+initErr.Error(), http.StatusInternalServerError)
		return
	}

	if origPath := r.Header.Get("x-matched-path"); origPath != "" {
		r.URL.Path = origPath
	} else if origURI := r.Header.Get("x-now-route-matches"); origURI != "" {
		parts := strings.SplitN(origURI, "?", 2)
		if len(parts) > 0 && parts[0] != "" {
			r.URL.Path = parts[0]
		}
	}

	srv.ServeHTTP(w, r)
}

package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/server"
)

func TestServerStaticAndBaseRoute(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	srv, err := server.NewServer(db)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	// 1. Static CSS check
	reqCSS := httptest.NewRequest("GET", "/static/css/style.css", nil)
	recCSS := httptest.NewRecorder()
	srv.ServeHTTP(recCSS, reqCSS)

	if recCSS.Code != http.StatusOK {
		t.Fatalf("expected /static/css/style.css to return 200, got %d", recCSS.Code)
	}
	contentType := recCSS.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/css") {
		t.Errorf("expected text/css content type for style.css, got %q", contentType)
	}
	if !strings.Contains(recCSS.Body.String(), "--bg") {
		t.Errorf("expected style.css to define --bg token")
	}

	// 2. Static HTMX JS check
	reqJS := httptest.NewRequest("GET", "/static/js/htmx.min.js", nil)
	recJS := httptest.NewRecorder()
	srv.ServeHTTP(recJS, reqJS)

	if recJS.Code != http.StatusOK {
		t.Fatalf("expected /static/js/htmx.min.js to return 200, got %d", recJS.Code)
	}
	jsContentType := recJS.Header().Get("Content-Type")
	if !strings.Contains(jsContentType, "javascript") {
		t.Errorf("expected javascript content type for htmx.min.js, got %q", jsContentType)
	}
	if !strings.Contains(recJS.Body.String(), "htmx") {
		t.Errorf("expected htmx script body to contain htmx identifier")
	}

	// 3. Health route
	reqHealth := httptest.NewRequest("GET", "/healthz", nil)
	recHealth := httptest.NewRecorder()
	srv.ServeHTTP(recHealth, reqHealth)

	if recHealth.Code != http.StatusOK {
		t.Fatalf("expected /healthz to return 200, got %d", recHealth.Code)
	}

	// 4. Template parsing check
	if srv.Templates() == nil {
		t.Fatalf("expected templates to be initialized and non-nil")
	}
	if srv.Templates().Lookup("base.html") == nil {
		t.Fatalf("expected base.html template to be loaded")
	}
}

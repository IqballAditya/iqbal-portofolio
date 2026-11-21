package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/internal/server"
)

func setupTestServer(t *testing.T) *server.Server {
	t.Helper()
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	if err := repository.SeedData(db); err != nil {
		t.Fatalf("SeedData failed: %v", err)
	}

	srv, err := server.NewServer(db)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	return srv
}

func TestGateway_Get(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/gateway", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from /gateway, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Who's Watching?") {
		t.Errorf("expected response to contain \"Who's Watching?\", got: %s", body)
	}
	for _, expected := range []string{"Recruiter", "Tech Lead", "Founder", "Guest"} {
		if !strings.Contains(body, expected) {
			t.Errorf("expected gateway to include persona %q", expected)
		}
	}
}

func TestGatewaySelect_Post(t *testing.T) {
	srv := setupTestServer(t)

	form := url.Values{}
	form.Set("persona", "tech_lead")
	req := httptest.NewRequest(http.MethodPost, "/gateway/select", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status 303 redirect, got %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "/" {
		t.Errorf("expected redirect to \"/\", got %q", location)
	}

	var foundCookie bool
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "iqbal_persona" && cookie.Value == "tech_lead" {
			foundCookie = true
			break
		}
	}
	if !foundCookie {
		t.Errorf("expected iqbal_persona cookie to be set to tech_lead, got cookies: %+v", rec.Result().Cookies())
	}
}

func TestHome_Get(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "iqbal_persona", Value: "tech_lead"})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from /, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Nexus AI Agent") {
		t.Errorf("expected home page to contain \"Nexus AI Agent\"")
	}

	expectedSubstrings := []string{
		"99% Match",
		"2026",
		"18+ YRS EXP",
		"Ultra HD 4K",
		"Lead Systems Architect",
		"Continue Reviewing Iqbal's Projects",
		"Top 10 Tech Stacks by Iqbal Today",
		"Proven Impact & Product Case Studies",
	}

	for _, str := range expectedSubstrings {
		if !strings.Contains(body, str) {
			t.Errorf("expected home page to contain %q", str)
		}
	}
}

func TestHome_MissingOrCorruptedCookieDefaultsToGuest(t *testing.T) {
	srv := setupTestServer(t)

	// Case 1: missing cookie
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	rec1 := httptest.NewRecorder()
	srv.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("expected status 200 for missing cookie, got %d", rec1.Code)
	}
	if !strings.Contains(rec1.Body.String(), "Nexus AI Agent") {
		t.Errorf("expected missing cookie request to render home page")
	}

	// Case 2: corrupted / invalid cookie
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.AddCookie(&http.Cookie{Name: "iqbal_persona", Value: "corrupted_hacker_string_123$!%"})
	rec2 := httptest.NewRecorder()
	srv.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected status 200 for corrupted cookie, got %d", rec2.Code)
	}
	if !strings.Contains(rec2.Body.String(), "Nexus AI Agent") {
		t.Errorf("expected corrupted cookie request to safely fallback to guest")
	}
}

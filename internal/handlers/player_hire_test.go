package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestPlayer_GetExistingProject(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/projects/nexus-ai/player", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	expectedSubstrings := []string{
		"Nexus AI Agent",
		"terminal-player",
		"player.js",
		"S1:E1",
		"Blueprint",
		"type WorkerPool struct",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(body, sub) {
			t.Errorf("expected player response to contain %q, but it didn't", sub)
		}
	}
}

func TestPlayer_GetNonExistentProject(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/projects/non-existent/player", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent project player, got %d", rec.Code)
	}
}

func TestShortlist_Toggle(t *testing.T) {
	srv := setupTestServer(t)

	// Step 1: POST without cookie -> sets iqbal_session cookie, adds item
	req1 := httptest.NewRequest(http.MethodPost, "/shortlist/nexus-ai", nil)
	rec1 := httptest.NewRecorder()
	srv.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("expected status 200 on first toggle, got %d", rec1.Code)
	}

	var sessionCookie *http.Cookie
	for _, c := range rec1.Result().Cookies() {
		if c.Name == "iqbal_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil || sessionCookie.Value == "" {
		t.Fatalf("expected iqbal_session cookie to be set")
	}

	// Verify project is shortlisted for this session
	isShort, err := srv.Shortlists.IsShortlisted(sessionCookie.Value, "nexus-ai")
	if err != nil {
		t.Fatalf("failed to check shortlist status: %v", err)
	}
	if !isShort {
		t.Errorf("expected nexus-ai to be shortlisted after first toggle")
	}

	// Step 2: POST again with the same cookie -> toggles off (removes item)
	req2 := httptest.NewRequest(http.MethodPost, "/shortlist/nexus-ai", nil)
	req2.AddCookie(sessionCookie)
	rec2 := httptest.NewRecorder()
	srv.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected status 200 on second toggle, got %d", rec2.Code)
	}

	isShort2, err := srv.Shortlists.IsShortlisted(sessionCookie.Value, "nexus-ai")
	if err != nil {
		t.Fatalf("failed to check shortlist status: %v", err)
	}
	if isShort2 {
		t.Errorf("expected nexus-ai to be removed from shortlist after second toggle")
	}
}

func TestShortlist_GetPage(t *testing.T) {
	srv := setupTestServer(t)

	// Pre-seed shortlist for session
	sessionID := "test-session-123"
	_, err := srv.Shortlists.ToggleShortlist(sessionID, "nexus-ai")
	if err != nil {
		t.Fatalf("failed to pre-seed shortlist: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/shortlist", nil)
	req.AddCookie(&http.Cookie{Name: "iqbal_session", Value: sessionID})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from /shortlist, got %d", rec.Code)
	}

	body := rec.Body.String()
	expectedSubstrings := []string{
		"My Shortlist",
		"Nexus AI Agent",
		"Choose Your Plan",
		"Full-Time Staff / Lead",
		"Advisory / Architecture Sprint",
		"Contract / Fractional",
		"Hire / Contact Iqbal",
		"hx-post=\"/inquiries\"",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(body, sub) {
			t.Errorf("expected shortlist page to contain %q, but it didn't", sub)
		}
	}
}

func TestInquiry_PostInvalid(t *testing.T) {
	srv := setupTestServer(t)

	form := url.Values{}
	form.Set("name", "Jane Doe")
	form.Set("email", "")
	form.Set("engagement_type", "staff")
	form.Set("message", "")

	req := httptest.NewRequest(http.MethodPost, "/inquiries", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for empty email/message, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(strings.ToLower(body), "error") && !strings.Contains(strings.ToLower(body), "valid") {
		t.Errorf("expected response to contain validation error details, got: %s", body)
	}
}

func TestInquiry_PostValid(t *testing.T) {
	srv := setupTestServer(t)

	form := url.Values{}
	form.Set("name", "Alice Smith")
	form.Set("email", "alice@netflix.com")
	form.Set("engagement_type", "staff")
	form.Set("message", "We are looking for a Principal Backend Engineer to lead our media routing tier.")

	req := httptest.NewRequest(http.MethodPost, "/inquiries", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 for valid inquiry, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Inquiry Received") && !strings.Contains(body, "Thank you") {
		t.Errorf("expected success notification, got: %s", body)
	}

	var count int
	err := srv.DB().QueryRow("SELECT COUNT(*) FROM inquiries WHERE email = ?", "alice@netflix.com").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query inquiries table: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 inquiry in database, found %d", count)
	}
}

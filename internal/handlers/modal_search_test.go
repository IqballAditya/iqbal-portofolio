package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestModal_GetExistingProject(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/projects/nexus-ai/modal", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from /projects/nexus-ai/modal, got %d", rec.Code)
	}

	body := rec.Body.String()

	expectedSubstrings := []string{
		"Nexus AI Agent",
		"Episodes & Milestones",
		"type WorkerPool struct",
		"99% Match",
		"/projects/nexus-ai/player",
		"hx-post=\"/shortlist/nexus-ai\"",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(body, sub) {
			t.Errorf("expected modal response to contain %q, but it didn't", sub)
		}
	}
}

func TestModal_GetNonExistentProject(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/projects/non-existent/modal", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent project modal, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Project Not Found") {
		t.Errorf("expected 404 response to contain \"Project Not Found\", got: %s", body)
	}
	if !strings.Contains(body, "modal-backdrop") {
		t.Errorf("expected 404 response to include modal-backdrop card, got: %s", body)
	}
}

func TestSearch_GetPage(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from /search, got %d", rec.Code)
	}

	body := rec.Body.String()

	expectedElements := []string{
		"hx-get=\"/api/search\"",
		"hx-trigger=\"keyup changed delay:250ms, search\"",
		"hx-target=\"#search-results\"",
		"All",
		"AI",
		"Backend",
		"Infra",
		"Fullstack",
		"Nexus AI Agent",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(body, elem) {
			t.Errorf("expected search page to contain %q", elem)
		}
	}
}

func TestSearchAPI_WithQuery(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=storage", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from /api/search, got %d", rec.Code)
	}

	body := rec.Body.String()

	if !strings.Contains(body, "Distributed Key-Value Engine") {
		t.Errorf("expected search API to return \"Distributed Key-Value Engine\" for query \"storage\"")
	}

	if strings.Contains(body, "Nexus AI Agent") {
		t.Errorf("expected search API results not to contain \"Nexus AI Agent\" for query \"storage\"")
	}
}

func TestSearchAPI_SQLInjectionSafe(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=%27+OR+1%3D1--", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 from SQL injection probe, got %d", rec.Code)
	}

	body := rec.Body.String()

	if !strings.Contains(body, "No projects match your search") {
		t.Errorf("expected empty search state for non-matching injection query, got: %s", body)
	}

	if strings.Contains(body, "Nexus AI Agent") && strings.Contains(body, "Distributed Key-Value Engine") {
		t.Errorf("potential SQL injection detected: all projects returned on ' OR 1=1-- query")
	}
}

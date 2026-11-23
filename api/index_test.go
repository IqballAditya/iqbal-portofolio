package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerHome(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/index?__path=", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Adith AI") && !strings.Contains(body, "Netflix") {
		t.Errorf("expected body to contain portfolio keywords, got: %s", body[:min(200, len(body))])
	}
}

func TestHandlerGateway(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/index?__path=gateway", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlerHealthz(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/index?__path=healthz", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d. Body: %s", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != "OK" {
		t.Errorf("expected OK, got %s", rec.Body.String())
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	expected := `{"status":"ok"}`

	if rec.Body.String() != expected {
		t.Fatalf(
			"expected body %s, got %s",
			expected,
			rec.Body.String(),
		)
	}
}

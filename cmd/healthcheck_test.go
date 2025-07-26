package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	app := &application{
		config: config{
			env: "testing",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	rr := httptest.NewRecorder()

	app.healthcheckHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	expectedBody := `"status":"available"`
	if !strings.Contains(rr.Body.String(), expectedBody) {
		t.Errorf("expected response body to contain %q, but got %q", expectedBody, rr.Body.String())
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"server/internal/config"
	"server/internal/helper"
	"testing"
)

func TestAuth(t *testing.T) {
	config.ApiKey = "test-api-key"
	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	middleware := Auth(mockHandler)
	middleware.ServeHTTP(rr, req)

	expectedStatus := http.StatusForbidden
	if status := rr.Code; status != expectedStatus {
		t.Errorf("wrong status code: got %v want %v", status, expectedStatus)
	}
	// Check the response body is first clientId
	expected := "invalid API Key\n"
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func mockHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		helper.LogErr(err)
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	middleware := Logger(mockHandler)
	middleware.ServeHTTP(rr, req)
	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("wrong status code: got %v want %v", status, expectedStatus)
	}
	// Check the response body is first clientId
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

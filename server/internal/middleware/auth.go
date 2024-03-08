package middleware

import (
    "fmt"
    "log"
    "net/http"
    "server/internal/config"
)

// Auth Middleware to force API key check
func Auth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get(config.HeaderApiKey) != config.ApiKey {
            err := fmt.Errorf("invalid API Key")
            log.Println("[request]:error", err)
            http.Error(w, err.Error(), http.StatusForbidden)
            return
        }
        next(w, r)
    }
}

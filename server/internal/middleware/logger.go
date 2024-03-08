package middleware

import (
    "log"
    "net/http"
    "net/url"
    "server/internal/helper"
)

// Logger Middleware to log every request
func Logger(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        decodedValue, err := url.QueryUnescape(r.RequestURI)
        if err != nil {
            helper.LogErr(err)
            return
        }
        log.Printf("[Request]:received %s from: %s", decodedValue, r.RemoteAddr)
        next(w, r)
    }
}

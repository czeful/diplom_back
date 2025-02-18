package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs each HTTP request with its method, URI, and duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %s", r.Method, r.RequestURI, duration)
	})
}

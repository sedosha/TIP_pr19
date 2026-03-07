package middleware

import (
    "log"
    "net/http"
    "time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        requestID := GetRequestID(r.Context())
        if requestID == "" {
            requestID = "no-request-id"
        }
        
        log.Printf("[%s] --> %s %s", requestID, r.Method, r.URL.Path)
        
        next.ServeHTTP(w, r)
        
        duration := time.Since(start)
        // Форматируем в миллисекундах
        log.Printf("[%s] <-- %s %s (%dms)", requestID, r.Method, r.URL.Path, duration.Milliseconds())
    })
}

package http

import (
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "tech-ip-sem2-grpc/shared/middleware"
)

func NewRouter(handler *AuthHandler) http.Handler {
    r := chi.NewRouter()
    
    r.Use(middleware.RequestIDMiddleware)
    r.Use(middleware.LoggingMiddleware)
    
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok","service":"auth"}`))
    })
    
    r.Route("/v1/auth", func(r chi.Router) {
        r.Post("/login", handler.Login)
        r.Get("/verify", handler.Verify)
    })
    
    return r
}

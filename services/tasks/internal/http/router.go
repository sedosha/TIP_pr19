package http

import (
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "tech-ip-sem2/services/tasks/internal/client/authclient"
    "tech-ip-sem2/shared/middleware"
)

func NewRouter(handler *TaskHandler, authClient *authclient.AuthClient) http.Handler {
    r := chi.NewRouter()
    
    r.Use(middleware.RequestIDMiddleware)
    r.Use(middleware.LoggingMiddleware)
    
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok","service":"tasks"}`))
    })
    
    r.Route("/v1/tasks", func(r chi.Router) {
        r.Use(authClient.AuthMiddleware)
        
        r.Post("/", handler.CreateTask)
        r.Get("/", handler.GetAllTasks)
        r.Get("/{id}", handler.GetTask)
        r.Patch("/{id}", handler.UpdateTask)
        r.Delete("/{id}", handler.DeleteTask)
    })
    
    return r
}

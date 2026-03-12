package http

import (
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "tech-ip-sem2-grpc/services/tasks/internal/grpcclient"
    "tech-ip-sem2-grpc/shared/middleware"
)

func NewRouter(handler *TaskHandler, authClient *grpcclient.AuthGRPCClient) http.Handler {
    r := chi.NewRouter()
    
    r.Use(middleware.RequestIDMiddleware)
    r.Use(middleware.LoggingMiddleware)
    
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok","service":"tasks"}`))
    })
    
    r.Route("/v1/tasks", func(r chi.Router) {
        r.Use(AuthGRPCMiddleware(authClient))
        
        r.Post("/", handler.CreateTask)
        r.Get("/", handler.GetAllTasks)
        r.Get("/{id}", handler.GetTask)
        r.Patch("/{id}", handler.UpdateTask)
        r.Delete("/{id}", handler.DeleteTask)
    })
    
    return r
}

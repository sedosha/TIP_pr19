package http

import (
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "go.uber.org/zap"
    
    "tech-ip-pz3-logging/services/tasks/internal/grpcclient"
    "tech-ip-pz3-logging/shared/middleware"
)

func NewRouter(handler *TaskHandler, authClient *grpcclient.AuthGRPCClient, log *zap.Logger) http.Handler {
    r := chi.NewRouter()
    
    r.Use(middleware.RequestIDMiddleware)
    
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        log.Debug("health check called")
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok","service":"tasks"}`))
    })
    
    r.Route("/v1/tasks", func(r chi.Router) {
        r.Use(AuthGRPCMiddleware(authClient, log))
        
        r.Post("/", handler.CreateTask)
        r.Get("/", handler.GetAllTasks)
        r.Get("/{id}", handler.GetTask)
        r.Patch("/{id}", handler.UpdateTask)
        r.Delete("/{id}", handler.DeleteTask)
    })
    
    return r
}

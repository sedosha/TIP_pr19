package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    
    "tech-ip-sem2/services/tasks/internal/client/authclient"
    httphandler "tech-ip-sem2/services/tasks/internal/http"
)

func main() {
    port := os.Getenv("TASKS_PORT")
    if port == "" {
        port = "8090"
    }
    
    authBaseURL := os.Getenv("AUTH_BASE_URL")
    if authBaseURL == "" {
        authBaseURL = "http://localhost:8089"
    }
    
    authClient := authclient.NewAuthClient(authBaseURL)
    
    handler := httphandler.NewTaskHandler()
    router := httphandler.NewRouter(handler, authClient)
    
    addr := fmt.Sprintf(":%s", port)
    log.Printf("Tasks service starting on port %s", port)
    
    if err := http.ListenAndServe(addr, router); err != nil {
        log.Fatalf("Failed to start tasks service: %v", err)
    }
}

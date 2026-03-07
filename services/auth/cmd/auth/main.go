package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    
    httphandler "tech-ip-sem2/services/auth/internal/http"
)

func main() {
    port := os.Getenv("AUTH_PORT")
    if port == "" {
        port = "8089"
    }
    
    handler := httphandler.NewAuthHandler()
    router := httphandler.NewRouter(handler)
    
    addr := fmt.Sprintf(":%s", port)
    log.Printf("Auth service starting on port %s", port)
    
    if err := http.ListenAndServe(addr, router); err != nil {
        log.Fatalf("Failed to start auth service: %v", err)
    }
}

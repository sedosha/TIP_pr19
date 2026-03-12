package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    
    "tech-ip-sem2-grpc/services/tasks/internal/grpcclient"
    httphandler "tech-ip-sem2-grpc/services/tasks/internal/http"
)

func main() {
    port := os.Getenv("TASKS_PORT")
    if port == "" {
        port = "8090"
    }
    
    grpcAddr := os.Getenv("AUTH_GRPC_ADDR")
    if grpcAddr == "" {
        grpcAddr = "localhost:50051"
        log.Printf("Using default gRPC address: %s", grpcAddr)
    }
    
    authClient, err := grpcclient.NewAuthGRPCClient(grpcAddr)
    if err != nil {
        log.Fatalf("Failed to create auth gRPC client: %v", err)
    }
    defer authClient.Close()
    
    handler := httphandler.NewTaskHandler()
    router := httphandler.NewRouter(handler, authClient)
    
    addr := fmt.Sprintf(":%s", port)
    log.Printf("Tasks service starting on port %s", port)
    log.Printf("Auth gRPC address: %s", grpcAddr)
    
    if err := http.ListenAndServe(addr, router); err != nil {
        log.Fatalf("Failed to start tasks service: %v", err)
    }
}

package main

import (
    "fmt"
    "net/http"
    "os"
    
    "go.uber.org/zap"
    
    "tech-ip-pz3-logging/services/tasks/internal/grpcclient"
    httphandler "tech-ip-pz3-logging/services/tasks/internal/http"
    "tech-ip-pz3-logging/shared/middleware"
    applogger "tech-ip-pz3-logging/pkg/logger"
)

func main() {
    zapLogger, err := applogger.New()
    if err != nil {
        panic(err)
    }
    defer zapLogger.Sync()
    
    port := os.Getenv("TASKS_PORT")
    if port == "" {
        port = "8090"
    }
    
    grpcAddr := os.Getenv("AUTH_GRPC_ADDR")
    if grpcAddr == "" {
        grpcAddr = "localhost:50051"
        zapLogger.Info("using default gRPC address", zap.String("addr", grpcAddr))
    }
    
    authClient, err := grpcclient.NewAuthGRPCClient(grpcAddr)
    if err != nil {
        zapLogger.Fatal("failed to create auth gRPC client", zap.Error(err))
    }
    defer authClient.Close()
    
    handler := httphandler.NewTaskHandler(zapLogger)
    router := httphandler.NewRouter(handler, authClient, zapLogger)
    
    rootHandler := middleware.ZapLoggingMiddleware(zapLogger, router)
    
    addr := fmt.Sprintf(":%s", port)
    zapLogger.Info("Tasks service starting",
        zap.String("port", port),
        zap.String("grpc_addr", grpcAddr),
    )
    
    if err := http.ListenAndServe(addr, rootHandler); err != nil {
        zapLogger.Fatal("failed to start tasks service", zap.Error(err))
    }
}

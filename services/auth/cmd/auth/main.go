package main

import (
    "log"
    "net"
    "os"
    
    "google.golang.org/grpc"
    
    pb "tech-ip-sem2-grpc/proto"
    authgrpc "tech-ip-sem2-grpc/services/auth/internal/grpc"
)

func main() {
    grpcPort := os.Getenv("AUTH_GRPC_PORT")
    if grpcPort == "" {
        grpcPort = "50051"
    }
    
    lis, err := net.Listen("tcp", ":"+grpcPort)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    s := grpc.NewServer()
    pb.RegisterAuthServiceServer(s, &authgrpc.AuthServer{})
    
    log.Printf("Auth gRPC server starting on port %s", grpcPort)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

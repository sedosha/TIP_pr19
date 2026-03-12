package grpc

import (
    "context"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    pb "tech-ip-sem2-grpc/proto"
)

type AuthServer struct {
    pb.UnimplementedAuthServiceServer
}

func (s *AuthServer) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
    token := req.GetToken()
    
    if token == "" {
        return nil, status.Errorf(codes.Unauthenticated, "empty token")
    }
    
    if token == "demo-token" {
        return &pb.VerifyResponse{
            Valid:   true,
            Subject: "user",
        }, nil
    }
    
    return nil, status.Errorf(codes.Unauthenticated, "invalid token")
}

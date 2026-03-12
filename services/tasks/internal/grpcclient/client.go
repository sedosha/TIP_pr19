package grpcclient

import (
    "context"
    "fmt"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    pb "tech-ip-sem2-grpc/proto"
    "tech-ip-sem2-grpc/shared/middleware"
)

type AuthGRPCClient struct {
    client pb.AuthServiceClient
    conn   *grpc.ClientConn
}

func NewAuthGRPCClient(addr string) (*AuthGRPCClient, error) {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return nil, fmt.Errorf("failed to connect to auth grpc: %w", err)
    }
    
    client := pb.NewAuthServiceClient(conn)
    return &AuthGRPCClient{
        client: client,
        conn:   conn,
    }, nil
}

func (c *AuthGRPCClient) Close() {
    if c.conn != nil {
        c.conn.Close()
    }
}

func (c *AuthGRPCClient) VerifyToken(ctx context.Context, token string) (bool, string, error) {
    requestID := middleware.GetRequestID(ctx)
    
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    fmt.Printf("[%s] Calling gRPC verify with token: %s\n", requestID, token)
    
    resp, err := c.client.Verify(ctx, &pb.VerifyRequest{Token: token})
    if err != nil {
        st, ok := status.FromError(err)
        if ok {
            fmt.Printf("[%s] gRPC error: %s (code: %v)\n", requestID, st.Message(), st.Code())
            
            if st.Code() == codes.Unauthenticated {
                return false, "", nil
            }
            if st.Code() == codes.DeadlineExceeded {
                return false, "", fmt.Errorf("auth service deadline exceeded")
            }
        }
        return false, "", fmt.Errorf("auth service error: %w", err)
    }
    
    fmt.Printf("[%s] gRPC response: valid=%v, subject=%s\n", 
        requestID, resp.Valid, resp.Subject)
    
    return resp.Valid, resp.Subject, nil
}

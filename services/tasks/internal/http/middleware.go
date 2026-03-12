package http

import (
    "context"
    "fmt"
    "net/http"
    
    "tech-ip-sem2-grpc/services/tasks/internal/grpcclient"
)

func AuthGRPCMiddleware(client *grpcclient.AuthGRPCClient) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
                return
            }
            
            var token string
            _, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
            if err != nil || token == "" {
                http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
                return
            }
            
            valid, subject, err := client.VerifyToken(r.Context(), token)
            if err != nil {
                http.Error(w, `{"error":"auth service unavailable"}`, http.StatusServiceUnavailable)
                return
            }
            
            if !valid {
                http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
                return
            }
            
            ctx := context.WithValue(r.Context(), "username", subject)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

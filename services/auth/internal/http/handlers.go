package http

import (
    "encoding/json"
    "net/http"
    "strings"
    
    "tech-ip-sem2/services/auth/internal/service"
    "tech-ip-sem2/shared/middleware"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())
    
    var req service.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request format"}`, http.StatusBadRequest)
        return
    }
    
    expectedPass, exists := service.ValidUsers[req.Username]
    if !exists || expectedPass != req.Password {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
        return
    }
    
    response := service.LoginResponse{
        AccessToken: "demo-token",
        TokenType:   "Bearer",
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", requestID)
    json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())
    
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(service.VerifyResponse{
            Valid: false,
            Error: "missing authorization header",
        })
        return
    }
    
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(service.VerifyResponse{
            Valid: false,
            Error: "invalid authorization format",
        })
        return
    }
    
    token := parts[1]
    
    if token == "demo-token" {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("X-Request-ID", requestID)
        json.NewEncoder(w).Encode(service.VerifyResponse{
            Valid:   true,
            Subject: "user",
        })
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusUnauthorized)
    json.NewEncoder(w).Encode(service.VerifyResponse{
        Valid: false,
        Error: "invalid token",
    })
}

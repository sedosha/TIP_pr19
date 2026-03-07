package authclient

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "tech-ip-sem2/shared/middleware"
)

type AuthClient struct {
    httpClient *http.Client
    baseURL    string
}

type VerifyResponse struct {
    Valid   bool   `json:"valid"`
    Subject string `json:"subject,omitempty"`
    Error   string `json:"error,omitempty"`
}

func NewAuthClient(baseURL string) *AuthClient {
    return &AuthClient{
        httpClient: &http.Client{
            Timeout: 2 * time.Second,
        },
        baseURL: baseURL,
    }
}

func (c *AuthClient) VerifyToken(ctx context.Context, token string) (*VerifyResponse, int, error) {
    if token == "" {
        return nil, http.StatusBadRequest, fmt.Errorf("empty token")
    }
    
    requestID := middleware.GetRequestID(ctx)
    
    url := fmt.Sprintf("%s/v1/auth/verify", c.baseURL)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("X-Request-ID", requestID)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, http.StatusServiceUnavailable, fmt.Errorf("auth service error: %w", err)
    }
    defer resp.Body.Close()
    
    var response VerifyResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, http.StatusInternalServerError, fmt.Errorf("failed to decode response: %w", err)
    }
    
    // Возвращаем статус ответа от Auth service
    return &response, resp.StatusCode, nil
}

func (c *AuthClient) AuthMiddleware(next http.Handler) http.Handler {
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
        
        // Получаем и ответ, и статус код от Auth
        resp, statusCode, err := c.VerifyToken(r.Context(), token)
        if err != nil {
            // Если Auth сервис недоступен - 503
            if statusCode == http.StatusServiceUnavailable {
                http.Error(w, `{"error":"auth service unavailable"}`, http.StatusServiceUnavailable)
                return
            }
            // Если Auth вернул ошибку (401, 403 и т.д.) - возвращаем тот же статус
            http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), statusCode)
            return
        }
        
        if !resp.Valid {
            // Auth сказал что токен невалидный - возвращаем 401
            http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
            return
        }
        
        ctx := context.WithValue(r.Context(), "username", resp.Subject)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

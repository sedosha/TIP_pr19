package httpx

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type ClientConfig struct {
    BaseURL    string
    Timeout    time.Duration
    MaxRetries int
}

func DefaultConfig(baseURL string) ClientConfig {
    return ClientConfig{
        BaseURL:    baseURL,
        Timeout:    3 * time.Second,
        MaxRetries: 1,
    }
}

type Client struct {
    httpClient *http.Client
    config     ClientConfig
}

func NewClient(config ClientConfig) *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: config.Timeout,
        },
        config: config,
    }
}

func (c *Client) Request(ctx context.Context, method, path string, body, response interface{}) error {
    url := c.config.BaseURL + path
    
    var bodyReader io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return fmt.Errorf("failed to marshal request body: %w", err)
        }
        bodyReader = bytes.NewReader(jsonBody)
    }
    
    req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    
    if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
        req.Header.Set("X-Request-ID", requestID)
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        return fmt.Errorf("request failed with status %d", resp.StatusCode)
    }
    
    if response != nil {
        if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
            return fmt.Errorf("failed to decode response: %w", err)
        }
    }
    
    return nil
}

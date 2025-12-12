package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"api-gateway/internal/resilience"
	"api-gateway/pkg/logger"

	"go.uber.org/zap"
)

type AuthClient struct {
	baseURL    string
	httpClient *http.Client
	cb         *resilience.CircuitBreaker
}

func NewAuthClient(baseURL string, cb *resilience.CircuitBreaker) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cb: cb,
	}
}

func (c *AuthClient) Register(ctx context.Context, body []byte) ([]byte, int, error) {
	return c.proxyRequest(ctx, "POST", "/api/auth/register", body, "")
}

func (c *AuthClient) Login(ctx context.Context, body []byte) ([]byte, int, error) {
	return c.proxyRequest(ctx, "POST", "/api/auth/login", body, "")
}

func (c *AuthClient) RefreshToken(ctx context.Context, body []byte) ([]byte, int, error) {
	return c.proxyRequest(ctx, "POST", "/api/auth/refresh", body, "")
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/auth/validate", c.baseURL)

	var result map[string]interface{}

	_, err := c.cb.Execute(func() (interface{}, error) {
		return nil, resilience.RetryWithBackoff(func() error {
			req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
			if err != nil {
				return err
			}

			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := c.httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("request failed: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("validation failed: status %d", resp.StatusCode)
			}

			return json.NewDecoder(resp.Body).Decode(&result)
		})
	})

	if err != nil {
		logger.Error("Token validation failed", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (c *AuthClient) ProxyRequest(ctx context.Context, method, path string, body []byte, token string) ([]byte, int, error) {
	return c.proxyRequest(ctx, method, path, body, token)
}

func (c *AuthClient) proxyRequest(ctx context.Context, method, path string, body []byte, token string) ([]byte, int, error) {
	url := c.baseURL + path

	var responseBody []byte
	var statusCode int

	_, err := c.cb.Execute(func() (interface{}, error) {
		return nil, resilience.RetryWithBackoff(func() error {
			var req *http.Request
			var err error

			if body != nil {
				req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
			} else {
				req, err = http.NewRequestWithContext(ctx, method, url, nil)
			}

			if err != nil {
				return err
			}

			req.Header.Set("Content-Type", "application/json")
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			resp, err := c.httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("request failed: %w", err)
			}
			defer resp.Body.Close()

			statusCode = resp.StatusCode
			responseBody, err = io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}

			if statusCode >= 500 {
				return fmt.Errorf("server error: %d", statusCode)
			}

			return nil
		})
	})

	if err != nil {
		logger.Error("Proxy request failed",
			logger.String("service", "auth-service"),
			logger.String("method", method),
			logger.String("path", path),
			zap.Error(err),
		)
		return nil, http.StatusServiceUnavailable, err
	}

	return responseBody, statusCode, nil
}

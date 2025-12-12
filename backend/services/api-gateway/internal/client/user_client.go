package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"api-gateway/internal/resilience"
	"api-gateway/pkg/logger"

	"go.uber.org/zap"
)

type UserClient struct {
	baseURL    string
	httpClient *http.Client
	cb         *resilience.CircuitBreaker
}

func NewUserClient(baseURL string, cb *resilience.CircuitBreaker) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cb: cb,
	}
}

func (c *UserClient) GetStudentProfile(ctx context.Context, token string) ([]byte, int, error) {
	return c.proxyRequest(ctx, "GET", "/api/users/students/me", nil, token)
}

func (c *UserClient) UpdateStudentProfile(ctx context.Context, body []byte, token string) ([]byte, int, error) {
	return c.proxyRequest(ctx, "PUT", "/api/users/students/me", body, token)
}

func (c *UserClient) DeleteStudentProfile(ctx context.Context, token string) ([]byte, int, error) {
	return c.proxyRequest(ctx, "DELETE", "/api/users/students/me", nil, token)
}

func (c *UserClient) GetEmployerProfile(ctx context.Context, token string) ([]byte, int, error) {
	return c.proxyRequest(ctx, "GET", "/api/users/employers/me", nil, token)
}

func (c *UserClient) UpdateEmployerProfile(ctx context.Context, body []byte, token string) ([]byte, int, error) {
	return c.proxyRequest(ctx, "PUT", "/api/users/employers/me", body, token)
}

func (c *UserClient) DeleteEmployerProfile(ctx context.Context, token string) ([]byte, int, error) {
	return c.proxyRequest(ctx, "DELETE", "/api/users/employers/me", nil, token)
}

func (c *UserClient) ProxyRequest(ctx context.Context, method, path string, body []byte, token string) ([]byte, int, error) {
	return c.proxyRequest(ctx, method, path, body, token)
}

func (c *UserClient) proxyRequest(ctx context.Context, method, path string, body []byte, token string) ([]byte, int, error) {
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
			logger.String("service", "user-service"),
			logger.String("method", method),
			logger.String("path", path),
			zap.Error(err),
		)
		return nil, http.StatusServiceUnavailable, err
	}

	return responseBody, statusCode, nil
}

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

type VacancyClient struct {
	baseURL    string
	httpClient *http.Client
	cb         *resilience.CircuitBreaker
}

func NewVacancyClient(baseURL string, cb *resilience.CircuitBreaker) *VacancyClient {
	return &VacancyClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cb: cb,
	}
}

func (c *VacancyClient) CreateVacancy(ctx context.Context, body []byte, token string, userID int, role string) ([]byte, int, error) {
	return c.proxyRequest(ctx, "POST", "/vacancies", body, token, userID, role)
}

func (c *VacancyClient) GetMyVacancies(ctx context.Context, token string, userID int, role string) ([]byte, int, error) {
	return c.proxyRequest(ctx, "GET", "/vacancies/me", nil, token, userID, role)
}

func (c *VacancyClient) GetVacanciesFeed(ctx context.Context, queryString string, token string) ([]byte, int, error) {
	path := "/vacancies/feed"
	if queryString != "" {
		path += "?" + queryString
	}
	return c.proxyRequest(ctx, "GET", path, nil, token, 0, "")
}

func (c *VacancyClient) GetVacancyByID(ctx context.Context, id string, token string) ([]byte, int, error) {
	path := fmt.Sprintf("/vacancies/%s", id)
	return c.proxyRequest(ctx, "GET", path, nil, token, 0, "")
}

func (c *VacancyClient) UpdateVacancy(ctx context.Context, id string, body []byte, token string, userID int, role string) ([]byte, int, error) {
	path := fmt.Sprintf("/vacancies/%s", id)
	return c.proxyRequest(ctx, "PUT", path, body, token, userID, role)
}

func (c *VacancyClient) DeleteVacancy(ctx context.Context, id string, token string, userID int, role string) ([]byte, int, error) {
	path := fmt.Sprintf("/vacancies/%s", id)
	return c.proxyRequest(ctx, "DELETE", path, nil, token, userID, role)
}

func (c *VacancyClient) PublishVacancy(ctx context.Context, id string, token string, userID int, role string) ([]byte, int, error) {
	path := fmt.Sprintf("/vacancies/%s/publish", id)
	return c.proxyRequest(ctx, "POST", path, nil, token, userID, role)
}

func (c *VacancyClient) SearchVacancies(ctx context.Context, queryString string, token string) ([]byte, int, error) {
	path := "/vacancies/search"
	if queryString != "" {
		path += "?" + queryString
	}
	return c.proxyRequest(ctx, "GET", path, nil, token, 0, "")
}

func (c *VacancyClient) ProxyRequest(ctx context.Context, method, path string, body []byte, token string, userID int, role string) ([]byte, int, error) {
	return c.proxyRequest(ctx, method, path, body, token, userID, role)
}

func (c *VacancyClient) proxyRequest(ctx context.Context, method, path string, body []byte, token string, userID int, role string) ([]byte, int, error) {
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
			if userID > 0 {
				req.Header.Set("X-User-Id", fmt.Sprintf("%d", userID))
			}
			if role != "" {
				req.Header.Set("X-User-Role", role)
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
			logger.String("service", "vacancy-service"),
			logger.String("method", method),
			logger.String("path", path),
			zap.Error(err),
		)
		return nil, http.StatusServiceUnavailable, err
	}

	return responseBody, statusCode, nil
}

package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserSyncClient struct {
	baseURL    string
	httpClient *http.Client
}

type createUserPayload struct {
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func NewUserSyncClient(baseURL string) *UserSyncClient {
	return &UserSyncClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// CreateUser forwards user creation to user-service (best-effort).
func (c *UserSyncClient) CreateUser(ctx context.Context, id int, email, password, role string) error {
	payload := createUserPayload{
		ID:       id,
		Email:    email,
		Password: password,
		Role:     role,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/users", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("user-service responded with status %d", resp.StatusCode)
	}

	return nil
}

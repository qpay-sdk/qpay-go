package qpay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const tokenBufferSeconds = 30

// Client is a thread-safe QPay V2 API client with automatic token management.
type Client struct {
	config *Config
	http   *http.Client
	mu     sync.Mutex

	accessToken      string
	refreshToken     string
	expiresAt        int64
	refreshExpiresAt int64
}

// NewClient creates a new QPay client with the given configuration.
func NewClient(cfg *Config) *Client {
	return &Client{
		config: cfg,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithHTTPClient creates a new QPay client with a custom http.Client.
func NewClientWithHTTPClient(cfg *Config, httpClient *http.Client) *Client {
	return &Client{
		config: cfg,
		http:   httpClient,
	}
}

func (c *Client) ensureToken(ctx context.Context) error {
	c.mu.Lock()
	now := time.Now().Unix()

	// Access token still valid
	if c.accessToken != "" && now < c.expiresAt-tokenBufferSeconds {
		c.mu.Unlock()
		return nil
	}

	// Determine strategy: refresh or full auth
	canRefresh := c.refreshToken != "" && now < c.refreshExpiresAt-tokenBufferSeconds
	refreshTok := c.refreshToken
	c.mu.Unlock()

	// Access token expired, try refresh
	if canRefresh {
		token, err := c.doRefreshTokenHTTP(ctx, refreshTok)
		if err == nil {
			c.mu.Lock()
			c.storeToken(token)
			c.mu.Unlock()
			return nil
		}
		// Refresh failed, fall through to get new token
	}

	// Both expired or no tokens, get new token
	token, err := c.getTokenRequest(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	c.mu.Lock()
	c.storeToken(token)
	c.mu.Unlock()
	return nil
}

// doRefreshTokenHTTP performs the HTTP call for token refresh without locking.
func (c *Client) doRefreshTokenHTTP(ctx context.Context, refreshTok string) (*TokenResponse, error) {
	url := c.config.BaseURL + "/v2/auth/refresh"
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+refreshTok)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		qErr := &Error{
			StatusCode: resp.StatusCode,
			RawBody:    string(respBody),
		}
		_ = json.Unmarshal(respBody, qErr)
		return nil, qErr
	}

	var token TokenResponse
	if err := json.Unmarshal(respBody, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (c *Client) storeToken(token *TokenResponse) {
	c.accessToken = token.AccessToken
	c.refreshToken = token.RefreshToken
	c.expiresAt = token.ExpiresIn
	c.refreshExpiresAt = token.RefreshExpiresIn
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	url := c.config.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		qErr := &Error{
			StatusCode: resp.StatusCode,
			RawBody:    string(respBody),
		}
		_ = json.Unmarshal(respBody, qErr)
		if qErr.Code == "" {
			qErr.Code = http.StatusText(resp.StatusCode)
		}
		if qErr.Message == "" {
			qErr.Message = string(respBody)
		}
		return qErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *Client) doBasicAuthRequest(ctx context.Context, method, path string, result interface{}) error {
	url := c.config.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.config.Username, c.config.Password)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		qErr := &Error{
			StatusCode: resp.StatusCode,
			RawBody:    string(respBody),
		}
		_ = json.Unmarshal(respBody, qErr)
		if qErr.Code == "" {
			qErr.Code = http.StatusText(resp.StatusCode)
		}
		if qErr.Message == "" {
			qErr.Message = string(respBody)
		}
		return qErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

package qpay

import "context"

// GetToken authenticates with QPay using Basic Auth and returns a new token pair.
// The token is automatically stored in the client for subsequent requests.
func (c *Client) GetToken(ctx context.Context) (*TokenResponse, error) {
	token, err := c.getTokenRequest(ctx)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.storeToken(token)
	c.mu.Unlock()

	return token, nil
}

// RefreshToken uses the current refresh token to obtain a new access token.
// The new token is automatically stored in the client for subsequent requests.
func (c *Client) RefreshToken(ctx context.Context) (*TokenResponse, error) {
	c.mu.Lock()
	refreshTok := c.refreshToken
	c.mu.Unlock()

	token, err := c.doRefreshTokenHTTP(ctx, refreshTok)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.storeToken(token)
	c.mu.Unlock()

	return token, nil
}

func (c *Client) getTokenRequest(ctx context.Context) (*TokenResponse, error) {
	var token TokenResponse
	if err := c.doBasicAuthRequest(ctx, "POST", "/v2/auth/token", &token); err != nil {
		return nil, err
	}
	return &token, nil
}

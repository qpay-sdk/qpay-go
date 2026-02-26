package qpay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetToken_Success(t *testing.T) {
	expectedToken := TokenResponse{
		AccessToken:      "access-abc",
		RefreshToken:     "refresh-xyz",
		ExpiresIn:        time.Now().Unix() + 3600,
		RefreshExpiresIn: time.Now().Unix() + 7200,
		TokenType:        "Bearer",
		Scope:            "openid",
		SessionState:     "session-123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/auth/token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Error("expected basic auth")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if user != "testuser" || pass != "testpass" {
			t.Errorf("unexpected credentials: %s:%s", user, pass)
		}

		json.NewEncoder(w).Encode(expectedToken)
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "testuser",
		Password: "testpass",
	}, server.Client())

	token, err := client.GetToken(context.Background())
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if token.AccessToken != expectedToken.AccessToken {
		t.Errorf("expected access token %q, got %q", expectedToken.AccessToken, token.AccessToken)
	}
	if token.RefreshToken != expectedToken.RefreshToken {
		t.Errorf("expected refresh token %q, got %q", expectedToken.RefreshToken, token.RefreshToken)
	}
	if token.TokenType != "Bearer" {
		t.Errorf("expected token type 'Bearer', got %q", token.TokenType)
	}

	// Verify token is stored in client
	if client.accessToken != expectedToken.AccessToken {
		t.Errorf("token not stored in client: got %q", client.accessToken)
	}
	if client.refreshToken != expectedToken.RefreshToken {
		t.Errorf("refresh token not stored in client: got %q", client.refreshToken)
	}
}

func TestGetToken_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "AUTHENTICATION_FAILED",
			"message": "Invalid username or password",
		})
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "wrong",
		Password: "creds",
	}, server.Client())

	token, err := client.GetToken(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if token != nil {
		t.Error("expected nil token on error")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T: %v", err, err)
	}
	if qErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", qErr.StatusCode)
	}
	if qErr.Code != "AUTHENTICATION_FAILED" {
		t.Errorf("expected error code 'AUTHENTICATION_FAILED', got %q", qErr.Code)
	}
}

func TestGetToken_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	_, err := client.GetToken(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRefreshToken_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/auth/refresh" {
			auth := r.Header.Get("Authorization")
			if auth != "Bearer old-refresh-token" {
				t.Errorf("expected Authorization 'Bearer old-refresh-token', got %q", auth)
			}
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "new-access",
				RefreshToken:     "new-refresh",
				ExpiresIn:        time.Now().Unix() + 3600,
				RefreshExpiresIn: time.Now().Unix() + 7200,
				TokenType:        "Bearer",
			})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	// Pre-set a refresh token
	client.refreshToken = "old-refresh-token"

	token, err := client.RefreshToken(context.Background())
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	if token.AccessToken != "new-access" {
		t.Errorf("expected access token 'new-access', got %q", token.AccessToken)
	}
	if token.RefreshToken != "new-refresh" {
		t.Errorf("expected refresh token 'new-refresh', got %q", token.RefreshToken)
	}

	// Verify tokens are stored
	if client.accessToken != "new-access" {
		t.Errorf("access token not updated in client: got %q", client.accessToken)
	}
	if client.refreshToken != "new-refresh" {
		t.Errorf("refresh token not updated in client: got %q", client.refreshToken)
	}
}

func TestRefreshToken_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "invalid_grant",
			"message": "Refresh token expired",
		})
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())
	client.refreshToken = "expired-token"

	token, err := client.RefreshToken(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if token != nil {
		t.Error("expected nil token on error")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T: %v", err, err)
	}
	if qErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", qErr.StatusCode)
	}
}

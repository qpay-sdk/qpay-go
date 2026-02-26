package qpay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cfg := &Config{
		BaseURL:     "https://api.qpay.mn",
		Username:    "user",
		Password:    "pass",
		InvoiceCode: "INV001",
		CallbackURL: "https://example.com/callback",
	}

	client := NewClient(cfg)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.config != cfg {
		t.Error("config not stored correctly")
	}
	if client.http == nil {
		t.Error("http client is nil")
	}
	if client.http.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", client.http.Timeout)
	}
	if client.accessToken != "" {
		t.Error("access token should be empty initially")
	}
	if client.refreshToken != "" {
		t.Error("refresh token should be empty initially")
	}
}

func TestNewClientWithHTTPClient(t *testing.T) {
	cfg := &Config{
		BaseURL:  "https://api.qpay.mn",
		Username: "user",
		Password: "pass",
	}
	custom := &http.Client{Timeout: 60 * time.Second}

	client := NewClientWithHTTPClient(cfg, custom)

	if client == nil {
		t.Fatal("NewClientWithHTTPClient returned nil")
	}
	if client.http != custom {
		t.Error("custom http client not stored correctly")
	}
	if client.http.Timeout != 60*time.Second {
		t.Errorf("expected timeout 60s, got %v", client.http.Timeout)
	}
}

func TestEnsureToken_FreshToken(t *testing.T) {
	// Mock server returns a token on POST /v2/auth/token
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/auth/token" && r.Method == "POST" {
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "access-123",
				RefreshToken:     "refresh-456",
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

	err := client.ensureToken(context.Background())
	if err != nil {
		t.Fatalf("ensureToken failed: %v", err)
	}
	if client.accessToken != "access-123" {
		t.Errorf("expected access token 'access-123', got %q", client.accessToken)
	}
	if client.refreshToken != "refresh-456" {
		t.Errorf("expected refresh token 'refresh-456', got %q", client.refreshToken)
	}
}

func TestEnsureToken_UsesExistingValidToken(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:      "access-123",
			RefreshToken:     "refresh-456",
			ExpiresIn:        time.Now().Unix() + 3600,
			RefreshExpiresIn: time.Now().Unix() + 7200,
		})
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	// First call fetches token
	if err := client.ensureToken(context.Background()); err != nil {
		t.Fatalf("first ensureToken failed: %v", err)
	}

	// Second call should not make a request
	if err := client.ensureToken(context.Background()); err != nil {
		t.Fatalf("second ensureToken failed: %v", err)
	}

	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("expected 1 server call, got %d", callCount)
	}
}

func TestEnsureToken_RefreshExpiredAccessToken(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		if r.URL.Path == "/v2/auth/token" {
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "access-old",
				RefreshToken:     "refresh-old",
				ExpiresIn:        time.Now().Unix() - 100, // already expired
				RefreshExpiresIn: time.Now().Unix() + 7200,
			})
			return
		}
		if r.URL.Path == "/v2/auth/refresh" {
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "access-new",
				RefreshToken:     "refresh-new",
				ExpiresIn:        time.Now().Unix() + 3600,
				RefreshExpiresIn: time.Now().Unix() + 7200,
			})
			return
		}
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	// First call gets initial token (which is already expired)
	if err := client.ensureToken(context.Background()); err != nil {
		t.Fatalf("first ensureToken failed: %v", err)
	}

	// Second call should attempt refresh since access token is expired
	if err := client.ensureToken(context.Background()); err != nil {
		t.Fatalf("second ensureToken failed: %v", err)
	}

	if client.accessToken != "access-new" {
		t.Errorf("expected refreshed access token 'access-new', got %q", client.accessToken)
	}
}

func TestEnsureToken_RefreshFails_FallsBackToFullAuth(t *testing.T) {
	callNum := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callNum++
		if r.URL.Path == "/v2/auth/token" {
			expiresIn := time.Now().Unix() + 3600
			if callNum == 1 {
				expiresIn = time.Now().Unix() - 100 // expired on first call
			}
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "access-fresh",
				RefreshToken:     "refresh-fresh",
				ExpiresIn:        expiresIn,
				RefreshExpiresIn: time.Now().Unix() + 7200,
			})
			return
		}
		if r.URL.Path == "/v2/auth/refresh" {
			// Refresh always fails
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "invalid_grant",
				"message": "Token is expired",
			})
			return
		}
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	// First call: get token (already expired)
	if err := client.ensureToken(context.Background()); err != nil {
		t.Fatalf("first ensureToken failed: %v", err)
	}

	// Second call: access token expired, refresh fails, falls back to full auth
	if err := client.ensureToken(context.Background()); err != nil {
		t.Fatalf("second ensureToken failed: %v", err)
	}

	if client.accessToken != "access-fresh" {
		t.Errorf("expected access token 'access-fresh', got %q", client.accessToken)
	}
}

func TestEnsureToken_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "server_error",
			"message": "internal server error",
		})
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	err := client.ensureToken(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDoRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/auth/token" {
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "test-token",
				RefreshToken:     "test-refresh",
				ExpiresIn:        time.Now().Unix() + 3600,
				RefreshExpiresIn: time.Now().Unix() + 7200,
			})
			return
		}
		if r.URL.Path == "/v2/test" {
			// Verify auth header
			auth := r.Header.Get("Authorization")
			if auth != "Bearer test-token" {
				t.Errorf("expected Authorization 'Bearer test-token', got %q", auth)
			}
			ct := r.Header.Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("expected Content-Type 'application/json', got %q", ct)
			}
			json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
			return
		}
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	var result map[string]string
	err := client.doRequest(context.Background(), "GET", "/v2/test", nil, &result)
	if err != nil {
		t.Fatalf("doRequest failed: %v", err)
	}
	if result["result"] != "ok" {
		t.Errorf("expected result 'ok', got %q", result["result"])
	}
}

func TestDoRequest_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/auth/token" {
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "test-token",
				RefreshToken:     "test-refresh",
				ExpiresIn:        time.Now().Unix() + 3600,
				RefreshExpiresIn: time.Now().Unix() + 7200,
			})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "INVOICE_NOTFOUND",
			"message": "Invoice not found",
		})
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	var result map[string]string
	err := client.doRequest(context.Background(), "GET", "/v2/test", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T: %v", err, err)
	}
	if qErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", qErr.StatusCode)
	}
	if qErr.Code != "INVOICE_NOTFOUND" {
		t.Errorf("expected code 'INVOICE_NOTFOUND', got %q", qErr.Code)
	}
}

func TestDoRequest_ServerError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/auth/token" {
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "test-token",
				RefreshToken:     "test-refresh",
				ExpiresIn:        time.Now().Unix() + 3600,
				RefreshExpiresIn: time.Now().Unix() + 7200,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	var result map[string]string
	err := client.doRequest(context.Background(), "GET", "/v2/test", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T: %v", err, err)
	}
	if qErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", qErr.StatusCode)
	}
}

func TestDoRequest_NilResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/auth/token" {
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "test-token",
				RefreshToken:     "test-refresh",
				ExpiresIn:        time.Now().Unix() + 3600,
				RefreshExpiresIn: time.Now().Unix() + 7200,
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	// Passing nil result should not panic
	err := client.doRequest(context.Background(), "DELETE", "/v2/test", nil, nil)
	if err != nil {
		t.Fatalf("doRequest with nil result failed: %v", err)
	}
}

func TestDoRequest_WithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/auth/token" {
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "test-token",
				RefreshToken:     "test-refresh",
				ExpiresIn:        time.Now().Unix() + 3600,
				RefreshExpiresIn: time.Now().Unix() + 7200,
			})
			return
		}
		// Decode the body to verify it was sent correctly
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if body["key"] != "value" {
			t.Errorf("expected body key 'value', got %v", body["key"])
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "user",
		Password: "pass",
	}, server.Client())

	reqBody := map[string]string{"key": "value"}
	var result map[string]string
	err := client.doRequest(context.Background(), "POST", "/v2/test", reqBody, &result)
	if err != nil {
		t.Fatalf("doRequest with body failed: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", result["status"])
	}
}

func TestDoBasicAuthRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Error("expected basic auth")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if user != "testuser" || pass != "testpass" {
			t.Errorf("unexpected credentials: %s:%s", user, pass)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "token-123",
			RefreshToken: "refresh-789",
		})
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "testuser",
		Password: "testpass",
	}, server.Client())

	var result TokenResponse
	err := client.doBasicAuthRequest(context.Background(), "POST", "/v2/auth/token", &result)
	if err != nil {
		t.Fatalf("doBasicAuthRequest failed: %v", err)
	}
	if result.AccessToken != "token-123" {
		t.Errorf("expected access token 'token-123', got %q", result.AccessToken)
	}
}

func TestDoBasicAuthRequest_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "AUTHENTICATION_FAILED",
			"message": "Invalid credentials",
		})
	}))
	defer server.Close()

	client := NewClientWithHTTPClient(&Config{
		BaseURL:  server.URL,
		Username: "wrong",
		Password: "creds",
	}, server.Client())

	var result TokenResponse
	err := client.doBasicAuthRequest(context.Background(), "POST", "/v2/auth/token", &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T: %v", err, err)
	}
	if qErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", qErr.StatusCode)
	}
	if qErr.Code != "AUTHENTICATION_FAILED" {
		t.Errorf("expected code 'AUTHENTICATION_FAILED', got %q", qErr.Code)
	}
}

// testHelper creates a mock server with token auth and a custom handler for the API path.
func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/auth/token" {
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:      "test-access-token",
				RefreshToken:     "test-refresh-token",
				ExpiresIn:        time.Now().Unix() + 3600,
				RefreshExpiresIn: time.Now().Unix() + 7200,
			})
			return
		}
		handler(w, r)
	}))

	client := NewClientWithHTTPClient(&Config{
		BaseURL:     server.URL,
		Username:    "user",
		Password:    "pass",
		InvoiceCode: "TEST_INVOICE",
		CallbackURL: "https://example.com/callback",
	}, server.Client())

	return client, server
}

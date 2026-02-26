package qpay

import (
	"os"
	"strings"
	"testing"
)

func TestLoadConfigFromEnv_Success(t *testing.T) {
	// Set all required env vars
	envVars := map[string]string{
		"QPAY_BASE_URL":     "https://merchant.qpay.mn",
		"QPAY_USERNAME":     "testuser",
		"QPAY_PASSWORD":     "testpass",
		"QPAY_INVOICE_CODE": "INV_CODE",
		"QPAY_CALLBACK_URL": "https://example.com/callback",
	}
	for k, v := range envVars {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadConfigFromEnv failed: %v", err)
	}

	if cfg.BaseURL != "https://merchant.qpay.mn" {
		t.Errorf("expected BaseURL 'https://merchant.qpay.mn', got %q", cfg.BaseURL)
	}
	if cfg.Username != "testuser" {
		t.Errorf("expected Username 'testuser', got %q", cfg.Username)
	}
	if cfg.Password != "testpass" {
		t.Errorf("expected Password 'testpass', got %q", cfg.Password)
	}
	if cfg.InvoiceCode != "INV_CODE" {
		t.Errorf("expected InvoiceCode 'INV_CODE', got %q", cfg.InvoiceCode)
	}
	if cfg.CallbackURL != "https://example.com/callback" {
		t.Errorf("expected CallbackURL 'https://example.com/callback', got %q", cfg.CallbackURL)
	}
}

func TestLoadConfigFromEnv_MissingBaseURL(t *testing.T) {
	envVars := map[string]string{
		"QPAY_USERNAME":     "testuser",
		"QPAY_PASSWORD":     "testpass",
		"QPAY_INVOICE_CODE": "INV_CODE",
		"QPAY_CALLBACK_URL": "https://example.com/callback",
	}
	for k, v := range envVars {
		os.Setenv(k, v)
	}
	os.Unsetenv("QPAY_BASE_URL")
	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected error for missing QPAY_BASE_URL, got nil")
	}
	if !strings.Contains(err.Error(), "QPAY_BASE_URL") {
		t.Errorf("error should mention QPAY_BASE_URL, got: %v", err)
	}
}

func TestLoadConfigFromEnv_MissingUsername(t *testing.T) {
	envVars := map[string]string{
		"QPAY_BASE_URL":     "https://merchant.qpay.mn",
		"QPAY_PASSWORD":     "testpass",
		"QPAY_INVOICE_CODE": "INV_CODE",
		"QPAY_CALLBACK_URL": "https://example.com/callback",
	}
	for k, v := range envVars {
		os.Setenv(k, v)
	}
	os.Unsetenv("QPAY_USERNAME")
	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected error for missing QPAY_USERNAME, got nil")
	}
	if !strings.Contains(err.Error(), "QPAY_USERNAME") {
		t.Errorf("error should mention QPAY_USERNAME, got: %v", err)
	}
}

func TestLoadConfigFromEnv_MissingPassword(t *testing.T) {
	envVars := map[string]string{
		"QPAY_BASE_URL":     "https://merchant.qpay.mn",
		"QPAY_USERNAME":     "testuser",
		"QPAY_INVOICE_CODE": "INV_CODE",
		"QPAY_CALLBACK_URL": "https://example.com/callback",
	}
	for k, v := range envVars {
		os.Setenv(k, v)
	}
	os.Unsetenv("QPAY_PASSWORD")
	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected error for missing QPAY_PASSWORD, got nil")
	}
	if !strings.Contains(err.Error(), "QPAY_PASSWORD") {
		t.Errorf("error should mention QPAY_PASSWORD, got: %v", err)
	}
}

func TestLoadConfigFromEnv_MissingInvoiceCode(t *testing.T) {
	envVars := map[string]string{
		"QPAY_BASE_URL":     "https://merchant.qpay.mn",
		"QPAY_USERNAME":     "testuser",
		"QPAY_PASSWORD":     "testpass",
		"QPAY_CALLBACK_URL": "https://example.com/callback",
	}
	for k, v := range envVars {
		os.Setenv(k, v)
	}
	os.Unsetenv("QPAY_INVOICE_CODE")
	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected error for missing QPAY_INVOICE_CODE, got nil")
	}
	if !strings.Contains(err.Error(), "QPAY_INVOICE_CODE") {
		t.Errorf("error should mention QPAY_INVOICE_CODE, got: %v", err)
	}
}

func TestLoadConfigFromEnv_MissingCallbackURL(t *testing.T) {
	envVars := map[string]string{
		"QPAY_BASE_URL":     "https://merchant.qpay.mn",
		"QPAY_USERNAME":     "testuser",
		"QPAY_PASSWORD":     "testpass",
		"QPAY_INVOICE_CODE": "INV_CODE",
	}
	for k, v := range envVars {
		os.Setenv(k, v)
	}
	os.Unsetenv("QPAY_CALLBACK_URL")
	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected error for missing QPAY_CALLBACK_URL, got nil")
	}
	if !strings.Contains(err.Error(), "QPAY_CALLBACK_URL") {
		t.Errorf("error should mention QPAY_CALLBACK_URL, got: %v", err)
	}
}

func TestLoadConfigFromEnv_AllMissing(t *testing.T) {
	// Unset all QPay env vars
	vars := []string{"QPAY_BASE_URL", "QPAY_USERNAME", "QPAY_PASSWORD", "QPAY_INVOICE_CODE", "QPAY_CALLBACK_URL"}
	for _, v := range vars {
		os.Unsetenv(v)
	}

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected error when all env vars are missing, got nil")
	}
}

func TestLoadConfigFromEnv_EmptyValue(t *testing.T) {
	envVars := map[string]string{
		"QPAY_BASE_URL":     "",
		"QPAY_USERNAME":     "testuser",
		"QPAY_PASSWORD":     "testpass",
		"QPAY_INVOICE_CODE": "INV_CODE",
		"QPAY_CALLBACK_URL": "https://example.com/callback",
	}
	for k, v := range envVars {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected error for empty QPAY_BASE_URL, got nil")
	}
	if !strings.Contains(err.Error(), "QPAY_BASE_URL") {
		t.Errorf("error should mention QPAY_BASE_URL, got: %v", err)
	}
}

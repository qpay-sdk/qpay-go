package qpay

import (
	"fmt"
	"os"
)

// Config holds the QPay API configuration.
type Config struct {
	BaseURL     string
	Username    string
	Password    string
	InvoiceCode string
	CallbackURL string
}

// LoadConfigFromEnv loads QPay configuration from environment variables.
//
// Required environment variables:
//   - QPAY_BASE_URL: QPay API base URL
//   - QPAY_USERNAME: QPay merchant username
//   - QPAY_PASSWORD: QPay merchant password
//   - QPAY_INVOICE_CODE: Default invoice code
//   - QPAY_CALLBACK_URL: Payment callback URL
func LoadConfigFromEnv() (*Config, error) {
	cfg := &Config{
		BaseURL:     os.Getenv("QPAY_BASE_URL"),
		Username:    os.Getenv("QPAY_USERNAME"),
		Password:    os.Getenv("QPAY_PASSWORD"),
		InvoiceCode: os.Getenv("QPAY_INVOICE_CODE"),
		CallbackURL: os.Getenv("QPAY_CALLBACK_URL"),
	}

	required := map[string]string{
		"QPAY_BASE_URL":      cfg.BaseURL,
		"QPAY_USERNAME":      cfg.Username,
		"QPAY_PASSWORD":      cfg.Password,
		"QPAY_INVOICE_CODE":  cfg.InvoiceCode,
		"QPAY_CALLBACK_URL":  cfg.CallbackURL,
	}

	for name, val := range required {
		if val == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", name)
		}
	}

	return cfg, nil
}

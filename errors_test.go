package qpay

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		contains []string
	}{
		{
			name: "standard error",
			err: &Error{
				StatusCode: http.StatusBadRequest,
				Code:       "INVOICE_NOTFOUND",
				Message:    "Invoice not found",
			},
			contains: []string{"qpay:", "INVOICE_NOTFOUND", "Invoice not found", "400"},
		},
		{
			name: "auth error",
			err: &Error{
				StatusCode: http.StatusUnauthorized,
				Code:       "AUTHENTICATION_FAILED",
				Message:    "Invalid credentials",
			},
			contains: []string{"qpay:", "AUTHENTICATION_FAILED", "Invalid credentials", "401"},
		},
		{
			name: "server error",
			err: &Error{
				StatusCode: http.StatusInternalServerError,
				Code:       "Internal Server Error",
				Message:    "Something went wrong",
			},
			contains: []string{"qpay:", "Internal Server Error", "Something went wrong", "500"},
		},
		{
			name: "empty code and message",
			err: &Error{
				StatusCode: http.StatusBadRequest,
				Code:       "",
				Message:    "",
			},
			contains: []string{"qpay:", "400"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			for _, s := range tt.contains {
				if !strings.Contains(msg, s) {
					t.Errorf("expected error message to contain %q, got %q", s, msg)
				}
			}
		})
	}
}

func TestError_ImplementsErrorInterface(t *testing.T) {
	var err error = &Error{
		StatusCode: 400,
		Code:       "TEST",
		Message:    "test",
	}

	// Verify it implements the error interface
	if err.Error() == "" {
		t.Error("Error() returned empty string")
	}
}

func TestError_Format(t *testing.T) {
	err := &Error{
		StatusCode: 404,
		Code:       "INVOICE_NOTFOUND",
		Message:    "Invoice not found",
	}

	expected := "qpay: INVOICE_NOTFOUND - Invoice not found (status 404)"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestIsQPayError_WithQPayError(t *testing.T) {
	original := &Error{
		StatusCode: http.StatusBadRequest,
		Code:       "INVOICE_NOTFOUND",
		Message:    "Invoice not found",
		RawBody:    `{"error":"INVOICE_NOTFOUND","message":"Invoice not found"}`,
	}

	qErr, ok := IsQPayError(original)
	if !ok {
		t.Fatal("expected IsQPayError to return true")
	}
	if qErr != original {
		t.Error("expected same error pointer")
	}
	if qErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", qErr.StatusCode)
	}
	if qErr.Code != "INVOICE_NOTFOUND" {
		t.Errorf("expected code 'INVOICE_NOTFOUND', got %q", qErr.Code)
	}
	if qErr.RawBody != `{"error":"INVOICE_NOTFOUND","message":"Invoice not found"}` {
		t.Errorf("unexpected raw body: %q", qErr.RawBody)
	}
}

func TestIsQPayError_WithStandardError(t *testing.T) {
	err := errors.New("some standard error")

	qErr, ok := IsQPayError(err)
	if ok {
		t.Error("expected IsQPayError to return false for standard error")
	}
	if qErr != nil {
		t.Error("expected nil error for non-QPay error")
	}
}

func TestIsQPayError_WithWrappedError(t *testing.T) {
	original := &Error{
		StatusCode: 400,
		Code:       "TEST",
		Message:    "test",
	}
	wrapped := fmt.Errorf("wrapped: %w", original)

	// IsQPayError uses type assertion, not errors.As, so wrapped errors won't match
	qErr, ok := IsQPayError(wrapped)
	if ok {
		t.Error("expected IsQPayError to return false for wrapped error (uses type assertion)")
	}
	if qErr != nil {
		t.Error("expected nil for wrapped error")
	}
}

func TestIsQPayError_WithNil(t *testing.T) {
	qErr, ok := IsQPayError(nil)
	if ok {
		t.Error("expected IsQPayError to return false for nil")
	}
	if qErr != nil {
		t.Error("expected nil error for nil input")
	}
}

func TestErrorConstants(t *testing.T) {
	// Verify a selection of error constants are defined correctly
	tests := map[string]string{
		"ErrAuthenticationFailed":      ErrAuthenticationFailed,
		"ErrInvoiceNotFound":           ErrInvoiceNotFound,
		"ErrPaymentNotFound":           ErrPaymentNotFound,
		"ErrInvalidAmount":             ErrInvalidAmount,
		"ErrPermissionDenied":          ErrPermissionDenied,
		"ErrPaymentAlreadyCanceled":    ErrPaymentAlreadyCanceled,
		"ErrPaymentNotPaid":            ErrPaymentNotPaid,
		"ErrInvoiceAlreadyCanceled":    ErrInvoiceAlreadyCanceled,
		"ErrEbarimtCancelNotSupported": ErrEbarimtCancelNotSupported,
		"ErrEbarimtNotRegistered":      ErrEbarimtNotRegistered,
	}

	expectedValues := map[string]string{
		"ErrAuthenticationFailed":      "AUTHENTICATION_FAILED",
		"ErrInvoiceNotFound":           "INVOICE_NOTFOUND",
		"ErrPaymentNotFound":           "PAYMENT_NOTFOUND",
		"ErrInvalidAmount":             "INVALID_AMOUNT",
		"ErrPermissionDenied":          "PERMISSION_DENIED",
		"ErrPaymentAlreadyCanceled":    "PAYMENT_ALREADY_CANCELED",
		"ErrPaymentNotPaid":            "PAYMENT_NOT_PAID",
		"ErrInvoiceAlreadyCanceled":    "INVOICE_ALREADY_CANCELED",
		"ErrEbarimtCancelNotSupported": "EBARIMT_CANCEL_NOTSUPPERDED",
		"ErrEbarimtNotRegistered":      "EBARIMT_NOT_REGISTERED",
	}

	for name, got := range tests {
		expected := expectedValues[name]
		if got != expected {
			t.Errorf("%s: expected %q, got %q", name, expected, got)
		}
	}
}

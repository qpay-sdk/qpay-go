package qpay

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateInvoice_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/invoice" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req CreateInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if req.InvoiceCode != "TEST_CODE" {
			t.Errorf("expected invoice code 'TEST_CODE', got %q", req.InvoiceCode)
		}
		if req.Amount != 50000 {
			t.Errorf("expected amount 50000, got %v", req.Amount)
		}

		json.NewEncoder(w).Encode(InvoiceResponse{
			InvoiceID:     "inv-123",
			QRText:        "qr-text-data",
			QRImage:       "base64-qr-image",
			QPay_ShortURL: "https://qpay.mn/s/abc",
			URLs: []Deeplink{
				{Name: "Khan Bank", Description: "Khan Bank app", Logo: "https://logo.png", Link: "khanbank://pay?q=abc"},
			},
		})
	})
	defer server.Close()

	resp, err := client.CreateInvoice(context.Background(), &CreateInvoiceRequest{
		InvoiceCode:         "TEST_CODE",
		SenderInvoiceNo:     "INV-001",
		InvoiceReceiverCode: "terminal",
		InvoiceDescription:  "Test invoice",
		Amount:              50000,
		CallbackURL:         "https://example.com/callback",
	})
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	if resp.InvoiceID != "inv-123" {
		t.Errorf("expected invoice ID 'inv-123', got %q", resp.InvoiceID)
	}
	if resp.QRText != "qr-text-data" {
		t.Errorf("expected QR text 'qr-text-data', got %q", resp.QRText)
	}
	if len(resp.URLs) != 1 {
		t.Fatalf("expected 1 URL, got %d", len(resp.URLs))
	}
	if resp.URLs[0].Name != "Khan Bank" {
		t.Errorf("expected URL name 'Khan Bank', got %q", resp.URLs[0].Name)
	}
}

func TestCreateInvoice_Error(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "INVOICE_CODE_INVALID",
			"message": "Invalid invoice code",
		})
	})
	defer server.Close()

	_, err := client.CreateInvoice(context.Background(), &CreateInvoiceRequest{
		InvoiceCode: "INVALID",
		Amount:      100,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T: %v", err, err)
	}
	if qErr.Code != "INVOICE_CODE_INVALID" {
		t.Errorf("expected error code 'INVOICE_CODE_INVALID', got %q", qErr.Code)
	}
}

func TestCreateSimpleInvoice_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/invoice" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req CreateSimpleInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Amount != 10000 {
			t.Errorf("expected amount 10000, got %v", req.Amount)
		}

		json.NewEncoder(w).Encode(InvoiceResponse{
			InvoiceID: "simple-inv-456",
			QRText:    "simple-qr",
		})
	})
	defer server.Close()

	resp, err := client.CreateSimpleInvoice(context.Background(), &CreateSimpleInvoiceRequest{
		InvoiceCode:         "TEST_CODE",
		SenderInvoiceNo:     "SINV-001",
		InvoiceReceiverCode: "terminal",
		InvoiceDescription:  "Simple test",
		Amount:              10000,
		CallbackURL:         "https://example.com/callback",
	})
	if err != nil {
		t.Fatalf("CreateSimpleInvoice failed: %v", err)
	}

	if resp.InvoiceID != "simple-inv-456" {
		t.Errorf("expected invoice ID 'simple-inv-456', got %q", resp.InvoiceID)
	}
}

func TestCreateSimpleInvoice_Error(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "INVALID_AMOUNT",
			"message": "Amount must be positive",
		})
	})
	defer server.Close()

	_, err := client.CreateSimpleInvoice(context.Background(), &CreateSimpleInvoiceRequest{
		InvoiceCode: "CODE",
		Amount:      -100,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "INVALID_AMOUNT" {
		t.Errorf("expected code 'INVALID_AMOUNT', got %q", qErr.Code)
	}
}

func TestCreateEbarimtInvoice_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/invoice" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req CreateEbarimtInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.TaxType != "1" {
			t.Errorf("expected tax type '1', got %q", req.TaxType)
		}
		if len(req.Lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(req.Lines))
		}

		json.NewEncoder(w).Encode(InvoiceResponse{
			InvoiceID: "ebarimt-inv-789",
			QRText:    "ebarimt-qr",
		})
	})
	defer server.Close()

	resp, err := client.CreateEbarimtInvoice(context.Background(), &CreateEbarimtInvoiceRequest{
		InvoiceCode:         "TEST_CODE",
		SenderInvoiceNo:     "EINV-001",
		InvoiceReceiverCode: "terminal",
		InvoiceDescription:  "Ebarimt invoice",
		TaxType:             "1",
		DistrictCode:        "34",
		CallbackURL:         "https://example.com/callback",
		Lines: []EbarimtInvoiceLine{
			{
				LineDescription: "Product A",
				LineQuantity:    "1",
				LineUnitPrice:   "50000",
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateEbarimtInvoice failed: %v", err)
	}

	if resp.InvoiceID != "ebarimt-inv-789" {
		t.Errorf("expected invoice ID 'ebarimt-inv-789', got %q", resp.InvoiceID)
	}
}

func TestCreateEbarimtInvoice_Error(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "INVOICE_LINE_REQUIRED",
			"message": "At least one invoice line is required",
		})
	})
	defer server.Close()

	_, err := client.CreateEbarimtInvoice(context.Background(), &CreateEbarimtInvoiceRequest{
		InvoiceCode: "CODE",
		TaxType:     "1",
		Lines:       nil,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "INVOICE_LINE_REQUIRED" {
		t.Errorf("expected code 'INVOICE_LINE_REQUIRED', got %q", qErr.Code)
	}
}

func TestCancelInvoice_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/invoice/inv-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.CancelInvoice(context.Background(), "inv-123")
	if err != nil {
		t.Fatalf("CancelInvoice failed: %v", err)
	}
}

func TestCancelInvoice_NotFound(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "INVOICE_NOTFOUND",
			"message": "Invoice not found",
		})
	})
	defer server.Close()

	err := client.CancelInvoice(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "INVOICE_NOTFOUND" {
		t.Errorf("expected code 'INVOICE_NOTFOUND', got %q", qErr.Code)
	}
}

func TestCancelInvoice_AlreadyCanceled(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "INVOICE_ALREADY_CANCELED",
			"message": "Invoice has already been canceled",
		})
	})
	defer server.Close()

	err := client.CancelInvoice(context.Background(), "canceled-inv")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "INVOICE_ALREADY_CANCELED" {
		t.Errorf("expected code 'INVOICE_ALREADY_CANCELED', got %q", qErr.Code)
	}
}

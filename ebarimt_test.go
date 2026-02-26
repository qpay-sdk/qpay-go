package qpay

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateEbarimt_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/ebarimt_v3/create" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req CreateEbarimtRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if req.PaymentID != "pay-123" {
			t.Errorf("expected payment ID 'pay-123', got %q", req.PaymentID)
		}
		if req.EbarimtReceiverType != "83" {
			t.Errorf("expected receiver type '83', got %q", req.EbarimtReceiverType)
		}

		json.NewEncoder(w).Encode(EbarimtResponse{
			ID:              "ebarimt-001",
			GPaymentID:      "pay-123",
			Amount:          "50000",
			VatAmount:       "5000",
			CityTaxAmount:   "500",
			EbarimtQRData:   "qr-data-here",
			EbarimtLottery:  "ABC123",
			BarimtStatus:    "CREATED",
			BarimtStatusDate: "2024-01-15T10:30:00",
			Status:          true,
		})
	})
	defer server.Close()

	resp, err := client.CreateEbarimt(context.Background(), &CreateEbarimtRequest{
		PaymentID:           "pay-123",
		EbarimtReceiverType: "83",
		DistrictCode:        "34",
	})
	if err != nil {
		t.Fatalf("CreateEbarimt failed: %v", err)
	}

	if resp.ID != "ebarimt-001" {
		t.Errorf("expected ID 'ebarimt-001', got %q", resp.ID)
	}
	if resp.GPaymentID != "pay-123" {
		t.Errorf("expected payment ID 'pay-123', got %q", resp.GPaymentID)
	}
	if resp.Amount != "50000" {
		t.Errorf("expected amount '50000', got %q", resp.Amount)
	}
	if resp.EbarimtLottery != "ABC123" {
		t.Errorf("expected lottery 'ABC123', got %q", resp.EbarimtLottery)
	}
	if !resp.Status {
		t.Error("expected status true")
	}
}

func TestCreateEbarimt_WithReceiver(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req CreateEbarimtRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if req.EbarimtReceiver != "1234567" {
			t.Errorf("expected receiver '1234567', got %q", req.EbarimtReceiver)
		}
		if req.EbarimtReceiverType != "80" {
			t.Errorf("expected receiver type '80', got %q", req.EbarimtReceiverType)
		}

		json.NewEncoder(w).Encode(EbarimtResponse{
			ID:                  "ebarimt-002",
			EbarimtReceiverType: "80",
			EbarimtReceiver:     "1234567",
			Status:              true,
		})
	})
	defer server.Close()

	resp, err := client.CreateEbarimt(context.Background(), &CreateEbarimtRequest{
		PaymentID:           "pay-456",
		EbarimtReceiverType: "80",
		EbarimtReceiver:     "1234567",
	})
	if err != nil {
		t.Fatalf("CreateEbarimt failed: %v", err)
	}

	if resp.EbarimtReceiverType != "80" {
		t.Errorf("expected receiver type '80', got %q", resp.EbarimtReceiverType)
	}
}

func TestCreateEbarimt_Error(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "EBARIMT_NOT_REGISTERED",
			"message": "Ebarimt not registered",
		})
	})
	defer server.Close()

	_, err := client.CreateEbarimt(context.Background(), &CreateEbarimtRequest{
		PaymentID:           "pay-invalid",
		EbarimtReceiverType: "83",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "EBARIMT_NOT_REGISTERED" {
		t.Errorf("expected code 'EBARIMT_NOT_REGISTERED', got %q", qErr.Code)
	}
	if qErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", qErr.StatusCode)
	}
}

func TestCreateEbarimt_ServerError(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})
	defer server.Close()

	_, err := client.CreateEbarimt(context.Background(), &CreateEbarimtRequest{
		PaymentID:           "pay-123",
		EbarimtReceiverType: "83",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCancelEbarimt_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/ebarimt_v3/pay-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		json.NewEncoder(w).Encode(EbarimtResponse{
			ID:               "ebarimt-001",
			GPaymentID:       "pay-123",
			BarimtStatus:     "CANCELED",
			BarimtStatusDate: "2024-01-15T12:00:00",
			Status:           false,
		})
	})
	defer server.Close()

	resp, err := client.CancelEbarimt(context.Background(), "pay-123")
	if err != nil {
		t.Fatalf("CancelEbarimt failed: %v", err)
	}

	if resp.ID != "ebarimt-001" {
		t.Errorf("expected ID 'ebarimt-001', got %q", resp.ID)
	}
	if resp.BarimtStatus != "CANCELED" {
		t.Errorf("expected status 'CANCELED', got %q", resp.BarimtStatus)
	}
	if resp.Status {
		t.Error("expected status false after cancellation")
	}
}

func TestCancelEbarimt_NotSupported(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "EBARIMT_CANCEL_NOTSUPPERDED",
			"message": "Ebarimt cancellation not supported",
		})
	})
	defer server.Close()

	_, err := client.CancelEbarimt(context.Background(), "pay-unsupported")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "EBARIMT_CANCEL_NOTSUPPERDED" {
		t.Errorf("expected code 'EBARIMT_CANCEL_NOTSUPPERDED', got %q", qErr.Code)
	}
}

func TestCancelEbarimt_ServerError(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "INTERNAL_ERROR",
			"message": "Server error",
		})
	})
	defer server.Close()

	_, err := client.CancelEbarimt(context.Background(), "pay-err")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", qErr.StatusCode)
	}
}

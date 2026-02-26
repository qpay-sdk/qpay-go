package qpay

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetPayment_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payment/pay-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		json.NewEncoder(w).Encode(PaymentDetail{
			PaymentID:       "pay-123",
			PaymentStatus:   "PAID",
			PaymentAmount:   "50000",
			PaymentCurrency: "MNT",
			PaymentWallet:   "qpay",
			PaymentDate:     "2024-01-15T10:30:00",
			ObjectType:      "INVOICE",
			ObjectID:        "inv-456",
		})
	})
	defer server.Close()

	payment, err := client.GetPayment(context.Background(), "pay-123")
	if err != nil {
		t.Fatalf("GetPayment failed: %v", err)
	}

	if payment.PaymentID != "pay-123" {
		t.Errorf("expected payment ID 'pay-123', got %q", payment.PaymentID)
	}
	if payment.PaymentStatus != "PAID" {
		t.Errorf("expected status 'PAID', got %q", payment.PaymentStatus)
	}
	if payment.PaymentAmount != "50000" {
		t.Errorf("expected amount '50000', got %q", payment.PaymentAmount)
	}
	if payment.PaymentCurrency != "MNT" {
		t.Errorf("expected currency 'MNT', got %q", payment.PaymentCurrency)
	}
}

func TestGetPayment_NotFound(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "PAYMENT_NOTFOUND",
			"message": "Payment not found",
		})
	})
	defer server.Close()

	_, err := client.GetPayment(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "PAYMENT_NOTFOUND" {
		t.Errorf("expected code 'PAYMENT_NOTFOUND', got %q", qErr.Code)
	}
}

func TestCheckPayment_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payment/check" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req PaymentCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if req.ObjectType != "INVOICE" {
			t.Errorf("expected object type 'INVOICE', got %q", req.ObjectType)
		}
		if req.ObjectID != "inv-123" {
			t.Errorf("expected object ID 'inv-123', got %q", req.ObjectID)
		}

		json.NewEncoder(w).Encode(PaymentCheckResponse{
			Count:      1,
			PaidAmount: 50000,
			Rows: []PaymentCheckRow{
				{
					PaymentID:     "pay-789",
					PaymentStatus: "PAID",
					PaymentAmount: "50000",
					PaymentWallet: "qpay",
				},
			},
		})
	})
	defer server.Close()

	resp, err := client.CheckPayment(context.Background(), &PaymentCheckRequest{
		ObjectType: "INVOICE",
		ObjectID:   "inv-123",
	})
	if err != nil {
		t.Fatalf("CheckPayment failed: %v", err)
	}

	if resp.Count != 1 {
		t.Errorf("expected count 1, got %d", resp.Count)
	}
	if resp.PaidAmount != 50000 {
		t.Errorf("expected paid amount 50000, got %v", resp.PaidAmount)
	}
	if len(resp.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(resp.Rows))
	}
	if resp.Rows[0].PaymentID != "pay-789" {
		t.Errorf("expected payment ID 'pay-789', got %q", resp.Rows[0].PaymentID)
	}
}

func TestCheckPayment_NoPayment(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(PaymentCheckResponse{
			Count:      0,
			PaidAmount: 0,
			Rows:       []PaymentCheckRow{},
		})
	})
	defer server.Close()

	resp, err := client.CheckPayment(context.Background(), &PaymentCheckRequest{
		ObjectType: "INVOICE",
		ObjectID:   "inv-unpaid",
	})
	if err != nil {
		t.Fatalf("CheckPayment failed: %v", err)
	}

	if resp.Count != 0 {
		t.Errorf("expected count 0, got %d", resp.Count)
	}
	if len(resp.Rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(resp.Rows))
	}
}

func TestCheckPayment_Error(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "INVALID_OBJECT_TYPE",
			"message": "Invalid object type",
		})
	})
	defer server.Close()

	_, err := client.CheckPayment(context.Background(), &PaymentCheckRequest{
		ObjectType: "INVALID",
		ObjectID:   "123",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "INVALID_OBJECT_TYPE" {
		t.Errorf("expected code 'INVALID_OBJECT_TYPE', got %q", qErr.Code)
	}
}

func TestListPayments_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payment/list" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req PaymentListRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		json.NewEncoder(w).Encode(PaymentListResponse{
			Count: 2,
			Rows: []PaymentListItem{
				{
					PaymentID:     "pay-001",
					PaymentDate:   "2024-01-10",
					PaymentStatus: "PAID",
					PaymentAmount: "10000",
				},
				{
					PaymentID:     "pay-002",
					PaymentDate:   "2024-01-11",
					PaymentStatus: "PAID",
					PaymentAmount: "20000",
				},
			},
		})
	})
	defer server.Close()

	resp, err := client.ListPayments(context.Background(), &PaymentListRequest{
		ObjectType: "INVOICE",
		ObjectID:   "inv-123",
		StartDate:  "2024-01-01",
		EndDate:    "2024-01-31",
		Offset:     Offset{PageNumber: 1, PageLimit: 10},
	})
	if err != nil {
		t.Fatalf("ListPayments failed: %v", err)
	}

	if resp.Count != 2 {
		t.Errorf("expected count 2, got %d", resp.Count)
	}
	if len(resp.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(resp.Rows))
	}
	if resp.Rows[0].PaymentID != "pay-001" {
		t.Errorf("expected first payment ID 'pay-001', got %q", resp.Rows[0].PaymentID)
	}
	if resp.Rows[1].PaymentAmount != "20000" {
		t.Errorf("expected second payment amount '20000', got %q", resp.Rows[1].PaymentAmount)
	}
}

func TestListPayments_Empty(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(PaymentListResponse{
			Count: 0,
			Rows:  []PaymentListItem{},
		})
	})
	defer server.Close()

	resp, err := client.ListPayments(context.Background(), &PaymentListRequest{
		ObjectType: "INVOICE",
		ObjectID:   "inv-999",
		StartDate:  "2024-01-01",
		EndDate:    "2024-01-31",
		Offset:     Offset{PageNumber: 1, PageLimit: 10},
	})
	if err != nil {
		t.Fatalf("ListPayments failed: %v", err)
	}

	if resp.Count != 0 {
		t.Errorf("expected count 0, got %d", resp.Count)
	}
}

func TestListPayments_ServerError(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "INTERNAL_ERROR",
			"message": "Something went wrong",
		})
	})
	defer server.Close()

	_, err := client.ListPayments(context.Background(), &PaymentListRequest{
		ObjectType: "INVOICE",
		ObjectID:   "inv-123",
	})
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

func TestCancelPayment_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payment/cancel/pay-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		var req PaymentCancelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if req.Note != "Cancel reason" {
			t.Errorf("expected note 'Cancel reason', got %q", req.Note)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.CancelPayment(context.Background(), "pay-123", &PaymentCancelRequest{
		CallbackURL: "https://example.com/callback",
		Note:        "Cancel reason",
	})
	if err != nil {
		t.Fatalf("CancelPayment failed: %v", err)
	}
}

func TestCancelPayment_NotFound(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "PAYMENT_NOTFOUND",
			"message": "Payment not found",
		})
	})
	defer server.Close()

	err := client.CancelPayment(context.Background(), "nonexistent", &PaymentCancelRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "PAYMENT_NOTFOUND" {
		t.Errorf("expected code 'PAYMENT_NOTFOUND', got %q", qErr.Code)
	}
}

func TestCancelPayment_AlreadyCanceled(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "PAYMENT_ALREADY_CANCELED",
			"message": "Payment already canceled",
		})
	})
	defer server.Close()

	err := client.CancelPayment(context.Background(), "pay-canceled", &PaymentCancelRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "PAYMENT_ALREADY_CANCELED" {
		t.Errorf("expected code 'PAYMENT_ALREADY_CANCELED', got %q", qErr.Code)
	}
}

func TestRefundPayment_Success(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payment/refund/pay-456" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		var req PaymentRefundRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.RefundPayment(context.Background(), "pay-456", &PaymentRefundRequest{
		CallbackURL: "https://example.com/callback",
		Note:        "Refund reason",
	})
	if err != nil {
		t.Fatalf("RefundPayment failed: %v", err)
	}
}

func TestRefundPayment_NotPaid(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "PAYMENT_NOT_PAID",
			"message": "Payment has not been paid",
		})
	})
	defer server.Close()

	err := client.RefundPayment(context.Background(), "unpaid", &PaymentRefundRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	qErr, ok := IsQPayError(err)
	if !ok {
		t.Fatalf("expected QPay error, got %T", err)
	}
	if qErr.Code != "PAYMENT_NOT_PAID" {
		t.Errorf("expected code 'PAYMENT_NOT_PAID', got %q", qErr.Code)
	}
}

func TestRefundPayment_ServerError(t *testing.T) {
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})
	defer server.Close()

	err := client.RefundPayment(context.Background(), "pay-err", &PaymentRefundRequest{})
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

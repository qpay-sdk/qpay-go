package qpay

import "context"

// GetPayment retrieves payment details by payment ID.
// GET /v2/payment/{id}
func (c *Client) GetPayment(ctx context.Context, paymentID string) (*PaymentDetail, error) {
	var resp PaymentDetail
	if err := c.doRequest(ctx, "GET", "/v2/payment/"+paymentID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CheckPayment checks if a payment has been made for an invoice.
// POST /v2/payment/check
func (c *Client) CheckPayment(ctx context.Context, req *PaymentCheckRequest) (*PaymentCheckResponse, error) {
	var resp PaymentCheckResponse
	if err := c.doRequest(ctx, "POST", "/v2/payment/check", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListPayments returns a list of payments matching the given criteria.
// POST /v2/payment/list
func (c *Client) ListPayments(ctx context.Context, req *PaymentListRequest) (*PaymentListResponse, error) {
	var resp PaymentListResponse
	if err := c.doRequest(ctx, "POST", "/v2/payment/list", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelPayment cancels a payment (card transactions only).
// DELETE /v2/payment/cancel/{id}
func (c *Client) CancelPayment(ctx context.Context, paymentID string, req *PaymentCancelRequest) error {
	return c.doRequest(ctx, "DELETE", "/v2/payment/cancel/"+paymentID, req, nil)
}

// RefundPayment refunds a payment (card transactions only).
// DELETE /v2/payment/refund/{id}
func (c *Client) RefundPayment(ctx context.Context, paymentID string, req *PaymentRefundRequest) error {
	return c.doRequest(ctx, "DELETE", "/v2/payment/refund/"+paymentID, req, nil)
}

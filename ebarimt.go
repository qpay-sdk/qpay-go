package qpay

import "context"

// CreateEbarimt creates an ebarimt (electronic tax receipt) for a payment.
// POST /v2/ebarimt_v3/create
func (c *Client) CreateEbarimt(ctx context.Context, req *CreateEbarimtRequest) (*EbarimtResponse, error) {
	var resp EbarimtResponse
	if err := c.doRequest(ctx, "POST", "/v2/ebarimt_v3/create", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelEbarimt cancels an ebarimt by payment ID.
// DELETE /v2/ebarimt_v3/{id}
func (c *Client) CancelEbarimt(ctx context.Context, paymentID string) (*EbarimtResponse, error) {
	var resp EbarimtResponse
	if err := c.doRequest(ctx, "DELETE", "/v2/ebarimt_v3/"+paymentID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

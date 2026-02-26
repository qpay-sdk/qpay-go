package qpay

import "context"

// CreateInvoice creates a detailed invoice with full options.
// POST /v2/invoice
func (c *Client) CreateInvoice(ctx context.Context, req *CreateInvoiceRequest) (*InvoiceResponse, error) {
	var resp InvoiceResponse
	if err := c.doRequest(ctx, "POST", "/v2/invoice", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateSimpleInvoice creates a simple invoice with minimal fields.
// POST /v2/invoice
func (c *Client) CreateSimpleInvoice(ctx context.Context, req *CreateSimpleInvoiceRequest) (*InvoiceResponse, error) {
	var resp InvoiceResponse
	if err := c.doRequest(ctx, "POST", "/v2/invoice", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateEbarimtInvoice creates an invoice with ebarimt (tax) information.
// POST /v2/invoice
func (c *Client) CreateEbarimtInvoice(ctx context.Context, req *CreateEbarimtInvoiceRequest) (*InvoiceResponse, error) {
	var resp InvoiceResponse
	if err := c.doRequest(ctx, "POST", "/v2/invoice", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelInvoice cancels an existing invoice by ID.
// DELETE /v2/invoice/{id}
func (c *Client) CancelInvoice(ctx context.Context, invoiceID string) error {
	return c.doRequest(ctx, "DELETE", "/v2/invoice/"+invoiceID, nil, nil)
}

# qpay-go

[![CI](https://github.com/qpay-sdk/qpay-go/actions/workflows/ci.yml/badge.svg)](https://github.com/qpay-sdk/qpay-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/qpay-sdk/qpay-go.svg)](https://pkg.go.dev/github.com/qpay-sdk/qpay-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Go SDK for the QPay V2 API. Provides a simple, type-safe client for creating invoices, managing payments, and handling ebarimt (electronic tax receipts) through QPay's merchant API.

## Installation

```bash
go get github.com/qpay-sdk/qpay-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	qpay "github.com/qpay-sdk/qpay-go"
)

func main() {
	// Load config from environment variables
	cfg, err := qpay.LoadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Create client
	client := qpay.NewClient(cfg)

	// Create a simple invoice
	invoice, err := client.CreateSimpleInvoice(context.Background(), &qpay.CreateSimpleInvoiceRequest{
		InvoiceCode:         cfg.InvoiceCode,
		SenderInvoiceNo:     "ORDER-001",
		InvoiceReceiverCode: "terminal",
		InvoiceDescription:  "Payment for Order #001",
		Amount:              50000,
		CallbackURL:         cfg.CallbackURL,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Invoice ID: %s\n", invoice.InvoiceID)
	fmt.Printf("QR Image:   %s\n", invoice.QRImage)
	fmt.Printf("Short URL:  %s\n", invoice.QPay_ShortURL)
}
```

## Configuration

### Environment Variables

Set these environment variables and use `LoadConfigFromEnv()`:

| Variable | Description |
|---|---|
| `QPAY_BASE_URL` | QPay API base URL (e.g. `https://merchant.qpay.mn`) |
| `QPAY_USERNAME` | QPay merchant username |
| `QPAY_PASSWORD` | QPay merchant password |
| `QPAY_INVOICE_CODE` | Default invoice code |
| `QPAY_CALLBACK_URL` | Payment callback URL |

```go
cfg, err := qpay.LoadConfigFromEnv()
if err != nil {
    log.Fatal(err)
}
client := qpay.NewClient(cfg)
```

### Manual Configuration

```go
client := qpay.NewClient(&qpay.Config{
    BaseURL:     "https://merchant.qpay.mn",
    Username:    "your_username",
    Password:    "your_password",
    InvoiceCode: "YOUR_INVOICE_CODE",
    CallbackURL: "https://yoursite.com/qpay/callback",
})
```

### Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}
client := qpay.NewClientWithHTTPClient(cfg, httpClient)
```

## Usage

### Authentication

The client handles authentication automatically. Tokens are obtained on first request and refreshed when they expire. You can also manage tokens manually:

```go
// Manually get a new token
token, err := client.GetToken(ctx)

// Manually refresh the token
token, err := client.RefreshToken(ctx)
```

### Create Invoice

**Simple invoice** with minimal fields:

```go
invoice, err := client.CreateSimpleInvoice(ctx, &qpay.CreateSimpleInvoiceRequest{
    InvoiceCode:         "YOUR_CODE",
    SenderInvoiceNo:     "ORDER-001",
    InvoiceReceiverCode: "terminal",
    InvoiceDescription:  "Payment for order",
    Amount:              50000,
    CallbackURL:         "https://yoursite.com/callback",
})
```

**Detailed invoice** with full options:

```go
invoice, err := client.CreateInvoice(ctx, &qpay.CreateInvoiceRequest{
    InvoiceCode:         "YOUR_CODE",
    SenderInvoiceNo:     "ORDER-002",
    InvoiceReceiverCode: "terminal",
    InvoiceDescription:  "Detailed payment",
    Amount:              150000,
    CallbackURL:         "https://yoursite.com/callback",
    SenderBranchData: &qpay.SenderBranchData{
        Name:  "Main Branch",
        Phone: "99001122",
    },
    Lines: []qpay.InvoiceLine{
        {
            LineDescription: "Product A",
            LineQuantity:    "2",
            LineUnitPrice:   "50000",
        },
        {
            LineDescription: "Product B",
            LineQuantity:    "1",
            LineUnitPrice:   "50000",
        },
    },
})
```

**Ebarimt invoice** with tax information:

```go
invoice, err := client.CreateEbarimtInvoice(ctx, &qpay.CreateEbarimtInvoiceRequest{
    InvoiceCode:         "YOUR_CODE",
    SenderInvoiceNo:     "ORDER-003",
    InvoiceReceiverCode: "terminal",
    InvoiceDescription:  "Tax invoice",
    TaxType:             "1",
    DistrictCode:        "34",
    CallbackURL:         "https://yoursite.com/callback",
    Lines: []qpay.EbarimtInvoiceLine{
        {
            TaxProductCode:  "1234567",
            LineDescription: "Taxable Product",
            LineQuantity:    "1",
            LineUnitPrice:   "100000",
        },
    },
})
```

### Cancel Invoice

```go
err := client.CancelInvoice(ctx, "invoice-id-here")
```

### Check Payment

```go
result, err := client.CheckPayment(ctx, &qpay.PaymentCheckRequest{
    ObjectType: "INVOICE",
    ObjectID:   "invoice-id-here",
    Offset: &qpay.Offset{
        PageNumber: 1,
        PageLimit:  10,
    },
})

if result.Count > 0 {
    fmt.Printf("Paid amount: %.0f\n", result.PaidAmount)
    for _, row := range result.Rows {
        fmt.Printf("Payment %s: %s (%s)\n", row.PaymentID, row.PaymentStatus, row.PaymentAmount)
    }
}
```

### Get Payment Details

```go
payment, err := client.GetPayment(ctx, "payment-id-here")

fmt.Printf("Status: %s\n", payment.PaymentStatus)
fmt.Printf("Amount: %s %s\n", payment.PaymentAmount, payment.PaymentCurrency)
fmt.Printf("Wallet: %s\n", payment.PaymentWallet)
```

### List Payments

```go
payments, err := client.ListPayments(ctx, &qpay.PaymentListRequest{
    ObjectType: "INVOICE",
    ObjectID:   "invoice-id-here",
    StartDate:  "2024-01-01",
    EndDate:    "2024-12-31",
    Offset: qpay.Offset{
        PageNumber: 1,
        PageLimit:  20,
    },
})

for _, p := range payments.Rows {
    fmt.Printf("%s - %s - %s MNT\n", p.PaymentDate, p.PaymentStatus, p.PaymentAmount)
}
```

### Cancel Payment

Cancel a card payment (card transactions only):

```go
err := client.CancelPayment(ctx, "payment-id-here", &qpay.PaymentCancelRequest{
    CallbackURL: "https://yoursite.com/callback",
    Note:        "Customer requested cancellation",
})
```

### Refund Payment

Refund a card payment (card transactions only):

```go
err := client.RefundPayment(ctx, "payment-id-here", &qpay.PaymentRefundRequest{
    CallbackURL: "https://yoursite.com/callback",
    Note:        "Product returned by customer",
})
```

### Create Ebarimt

Create an electronic tax receipt for a completed payment:

```go
// Individual (citizen) receipt
ebarimt, err := client.CreateEbarimt(ctx, &qpay.CreateEbarimtRequest{
    PaymentID:           "payment-id-here",
    EbarimtReceiverType: "83",
    DistrictCode:        "34",
})

// Organization receipt
ebarimt, err := client.CreateEbarimt(ctx, &qpay.CreateEbarimtRequest{
    PaymentID:           "payment-id-here",
    EbarimtReceiverType: "80",
    EbarimtReceiver:     "1234567",  // TIN number
    DistrictCode:        "34",
})

fmt.Printf("Lottery: %s\n", ebarimt.EbarimtLottery)
fmt.Printf("QR Data: %s\n", ebarimt.EbarimtQRData)
```

### Cancel Ebarimt

```go
ebarimt, err := client.CancelEbarimt(ctx, "payment-id-here")
```

## Error Handling

All API errors are returned as `*qpay.Error` which includes the HTTP status code, QPay error code, and message.

```go
invoice, err := client.CreateSimpleInvoice(ctx, req)
if err != nil {
    // Check if it's a QPay API error
    if qErr, ok := qpay.IsQPayError(err); ok {
        fmt.Printf("QPay error: %s (status %d)\n", qErr.Code, qErr.StatusCode)
        fmt.Printf("Message: %s\n", qErr.Message)
        fmt.Printf("Raw response: %s\n", qErr.RawBody)

        // Handle specific error codes
        switch qErr.Code {
        case qpay.ErrInvoiceCodeInvalid:
            // Handle invalid invoice code
        case qpay.ErrInvalidAmount:
            // Handle invalid amount
        case qpay.ErrAuthenticationFailed:
            // Handle auth failure
        }
    } else {
        // Network error or other non-API error
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Error Code Constants

The SDK provides constants for all QPay error codes. Some commonly used ones:

| Constant | Value | Description |
|---|---|---|
| `ErrAuthenticationFailed` | `AUTHENTICATION_FAILED` | Invalid credentials |
| `ErrInvoiceNotFound` | `INVOICE_NOTFOUND` | Invoice does not exist |
| `ErrInvoiceAlreadyCanceled` | `INVOICE_ALREADY_CANCELED` | Invoice was already canceled |
| `ErrInvoicePaid` | `INVOICE_PAID` | Invoice has already been paid |
| `ErrInvoiceCodeInvalid` | `INVOICE_CODE_INVALID` | Invalid invoice code |
| `ErrInvalidAmount` | `INVALID_AMOUNT` | Invalid payment amount |
| `ErrPaymentNotFound` | `PAYMENT_NOTFOUND` | Payment does not exist |
| `ErrPaymentNotPaid` | `PAYMENT_NOT_PAID` | Payment has not been completed |
| `ErrPaymentAlreadyCanceled` | `PAYMENT_ALREADY_CANCELED` | Payment was already canceled |
| `ErrEbarimtNotRegistered` | `EBARIMT_NOT_REGISTERED` | Ebarimt not registered |
| `ErrPermissionDenied` | `PERMISSION_DENIED` | Insufficient permissions |

See `errors.go` for the complete list of error constants.

## API Reference

| Method | Description | Returns |
|---|---|---|
| `NewClient(cfg)` | Create client with default HTTP settings | `*Client` |
| `NewClientWithHTTPClient(cfg, http)` | Create client with custom HTTP client | `*Client` |
| `GetToken(ctx)` | Authenticate and get token | `*TokenResponse, error` |
| `RefreshToken(ctx)` | Refresh access token | `*TokenResponse, error` |
| `CreateInvoice(ctx, req)` | Create detailed invoice | `*InvoiceResponse, error` |
| `CreateSimpleInvoice(ctx, req)` | Create simple invoice | `*InvoiceResponse, error` |
| `CreateEbarimtInvoice(ctx, req)` | Create invoice with ebarimt | `*InvoiceResponse, error` |
| `CancelInvoice(ctx, id)` | Cancel invoice by ID | `error` |
| `GetPayment(ctx, id)` | Get payment details | `*PaymentDetail, error` |
| `CheckPayment(ctx, req)` | Check payment status | `*PaymentCheckResponse, error` |
| `ListPayments(ctx, req)` | List payments | `*PaymentListResponse, error` |
| `CancelPayment(ctx, id, req)` | Cancel card payment | `error` |
| `RefundPayment(ctx, id, req)` | Refund card payment | `error` |
| `CreateEbarimt(ctx, req)` | Create ebarimt receipt | `*EbarimtResponse, error` |
| `CancelEbarimt(ctx, id)` | Cancel ebarimt | `*EbarimtResponse, error` |
| `LoadConfigFromEnv()` | Load config from env vars | `*Config, error` |
| `IsQPayError(err)` | Check if error is QPay error | `*Error, bool` |

## License

MIT License. See [LICENSE](LICENSE) for details.

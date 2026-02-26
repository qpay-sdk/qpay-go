package qpay

import "fmt"

// Error represents a QPay API error response.
type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"error"`
	Message    string `json:"message"`
	RawBody    string `json:"-"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("qpay: %s - %s (status %d)", e.Code, e.Message, e.StatusCode)
}

// IsQPayError checks if an error is a QPay API error and returns it.
func IsQPayError(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	if qErr, ok := err.(*Error); ok {
		return qErr, true
	}
	return nil, false
}

// QPay error code constants.
const (
	ErrAccountBankDuplicated          = "ACCOUNT_BANK_DUPLICATED"
	ErrAccountSelectionInvalid        = "ACCOUNT_SELECTION_INVALID"
	ErrAuthenticationFailed           = "AUTHENTICATION_FAILED"
	ErrBankAccountNotFound            = "BANK_ACCOUNT_NOTFOUND"
	ErrBankMCCAlreadyAdded            = "BANK_MCC_ALREADY_ADDED"
	ErrBankMCCNotFound                = "BANK_MCC_NOT_FOUND"
	ErrCardTerminalNotFound           = "CARD_TERMINAL_NOTFOUND"
	ErrClientNotFound                 = "CLIENT_NOTFOUND"
	ErrClientUsernameDuplicated       = "CLIENT_USERNAME_DUPLICATED"
	ErrCustomerDuplicate              = "CUSTOMER_DUPLICATE"
	ErrCustomerNotFound               = "CUSTOMER_NOTFOUND"
	ErrCustomerRegisterInvalid        = "CUSTOMER_REGISTER_INVALID"
	ErrEbarimtCancelNotSupported      = "EBARIMT_CANCEL_NOTSUPPERDED"
	ErrEbarimtNotRegistered           = "EBARIMT_NOT_REGISTERED"
	ErrEbarimtQRCodeInvalid           = "EBARIMT_QR_CODE_INVALID"
	ErrInformNotFound                 = "INFORM_NOTFOUND"
	ErrInputCodeRegistered            = "INPUT_CODE_REGISTERED"
	ErrInputNotFound                  = "INPUT_NOTFOUND"
	ErrInvalidAmount                  = "INVALID_AMOUNT"
	ErrInvalidObjectType              = "INVALID_OBJECT_TYPE"
	ErrInvoiceAlreadyCanceled         = "INVOICE_ALREADY_CANCELED"
	ErrInvoiceCodeInvalid             = "INVOICE_CODE_INVALID"
	ErrInvoiceCodeRegistered          = "INVOICE_CODE_REGISTERED"
	ErrInvoiceLineRequired            = "INVOICE_LINE_REQUIRED"
	ErrInvoiceNotFound                = "INVOICE_NOTFOUND"
	ErrInvoicePaid                    = "INVOICE_PAID"
	ErrInvoiceReceiverDataAddrReq     = "INVOICE_RECEIVER_DATA_ADDRESS_REQUIRED"
	ErrInvoiceReceiverDataEmailReq    = "INVOICE_RECEIVER_DATA_EMAIL_REQUIRED"
	ErrInvoiceReceiverDataPhoneReq    = "INVOICE_RECEIVER_DATA_PHONE_REQUIRED"
	ErrInvoiceReceiverDataRequired    = "INVOICE_RECEIVER_DATA_REQUIRED"
	ErrMaxAmountErr                   = "MAX_AMOUNT_ERR"
	ErrMCCNotFound                    = "MCC_NOTFOUND"
	ErrMerchantAlreadyRegistered      = "MERCHANT_ALREADY_REGISTERED"
	ErrMerchantInactive               = "MERCHANT_INACTIVE"
	ErrMerchantNotFound               = "MERCHANT_NOTFOUND"
	ErrMinAmountErr                   = "MIN_AMOUNT_ERR"
	ErrNoCredentials                  = "NO_CREDENDIALS"
	ErrObjectDataError                = "OBJECT_DATA_ERROR"
	ErrP2PTerminalNotFound            = "P2P_TERMINAL_NOTFOUND"
	ErrPaymentAlreadyCanceled         = "PAYMENT_ALREADY_CANCELED"
	ErrPaymentNotPaid                 = "PAYMENT_NOT_PAID"
	ErrPaymentNotFound                = "PAYMENT_NOTFOUND"
	ErrPermissionDenied               = "PERMISSION_DENIED"
	ErrQRAccountInactive              = "QRACCOUNT_INACTIVE"
	ErrQRAccountNotFound              = "QRACCOUNT_NOTFOUND"
	ErrQRCodeNotFound                 = "QRCODE_NOTFOUND"
	ErrQRCodeUsed                     = "QRCODE_USED"
	ErrSenderBranchDataRequired       = "SENDER_BRANCH_DATA_REQUIRED"
	ErrTaxLineRequired                = "TAX_LINE_REQUIRED"
	ErrTaxProductCodeRequired         = "TAX_PRODUCT_CODE_REQUIRED"
	ErrTransactionNotApproved         = "TRANSACTION_NOT_APPROVED"
	ErrTransactionRequired            = "TRANSACTION_REQUIRED"
)

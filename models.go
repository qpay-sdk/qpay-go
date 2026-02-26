package qpay

// --- Auth ---

// TokenResponse represents the QPay authentication token response.
type TokenResponse struct {
	TokenType        string `json:"token_type"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	Scope            string `json:"scope"`
	NotBeforePolicy  string `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
}

// --- Common nested types ---

// Address represents a physical address.
type Address struct {
	City      string `json:"city,omitempty"`
	District  string `json:"district,omitempty"`
	Street    string `json:"street,omitempty"`
	Building  string `json:"building,omitempty"`
	Address   string `json:"address,omitempty"`
	Zipcode   string `json:"zipcode,omitempty"`
	Longitude string `json:"longitude,omitempty"`
	Latitude  string `json:"latitude,omitempty"`
}

// SenderBranchData represents the sender branch information.
type SenderBranchData struct {
	Register string   `json:"register,omitempty"`
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Phone    string   `json:"phone,omitempty"`
	Address  *Address `json:"address,omitempty"`
}

// SenderStaffData represents the sender staff information.
type SenderStaffData struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// InvoiceReceiverData represents the invoice receiver information.
type InvoiceReceiverData struct {
	Register string   `json:"register,omitempty"`
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Phone    string   `json:"phone,omitempty"`
	Address  *Address `json:"address,omitempty"`
}

// Account represents a bank account used in transactions.
type Account struct {
	AccountBankCode string `json:"account_bank_code"`
	AccountNumber   string `json:"account_number"`
	IBANNumber      string `json:"iban_number"`
	AccountName     string `json:"account_name"`
	AccountCurrency string `json:"account_currency"`
	IsDefault       bool   `json:"is_default"`
}

// Transaction represents a payment transaction.
type Transaction struct {
	Description string    `json:"description"`
	Amount      string    `json:"amount"`
	Accounts    []Account `json:"accounts,omitempty"`
}

// InvoiceLine represents a line item in an invoice.
type InvoiceLine struct {
	TaxProductCode  string     `json:"tax_product_code,omitempty"`
	LineDescription string     `json:"line_description"`
	LineQuantity    string     `json:"line_quantity"`
	LineUnitPrice   string     `json:"line_unit_price"`
	Note            string     `json:"note,omitempty"`
	Discounts       []TaxEntry `json:"discounts,omitempty"`
	Surcharges      []TaxEntry `json:"surcharges,omitempty"`
	Taxes           []TaxEntry `json:"taxes,omitempty"`
}

// EbarimtInvoiceLine represents a line item in an ebarimt invoice.
type EbarimtInvoiceLine struct {
	TaxProductCode     string     `json:"tax_product_code,omitempty"`
	LineDescription    string     `json:"line_description"`
	Barcode            string     `json:"barcode,omitempty"`
	LineQuantity       string     `json:"line_quantity"`
	LineUnitPrice      string     `json:"line_unit_price"`
	Note               string     `json:"note,omitempty"`
	ClassificationCode string     `json:"classification_code,omitempty"`
	Taxes              []TaxEntry `json:"taxes,omitempty"`
}

// TaxEntry represents a tax, discount, or surcharge entry.
type TaxEntry struct {
	TaxCode       string  `json:"tax_code,omitempty"`
	DiscountCode  string  `json:"discount_code,omitempty"`
	SurchargeCode string  `json:"surcharge_code,omitempty"`
	Description   string  `json:"description"`
	Amount        float64 `json:"amount"`
	Note          string  `json:"note,omitempty"`
}

// Deeplink represents a payment deeplink for a bank or wallet app.
type Deeplink struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Link        string `json:"link"`
}

// --- Invoice ---

// CreateInvoiceRequest is the request body for creating a detailed invoice.
type CreateInvoiceRequest struct {
	InvoiceCode          string                `json:"invoice_code"`
	SenderInvoiceNo      string                `json:"sender_invoice_no"`
	SenderBranchCode     string                `json:"sender_branch_code,omitempty"`
	SenderBranchData     *SenderBranchData     `json:"sender_branch_data,omitempty"`
	SenderStaffData      *SenderStaffData      `json:"sender_staff_data,omitempty"`
	SenderStaffCode      string                `json:"sender_staff_code,omitempty"`
	InvoiceReceiverCode  string                `json:"invoice_receiver_code"`
	InvoiceReceiverData  *InvoiceReceiverData  `json:"invoice_receiver_data,omitempty"`
	InvoiceDescription   string                `json:"invoice_description"`
	EnableExpiry         *string               `json:"enable_expiry,omitempty"`
	AllowPartial         *bool                 `json:"allow_partial,omitempty"`
	MinimumAmount        *float64              `json:"minimum_amount,omitempty"`
	AllowExceed          *bool                 `json:"allow_exceed,omitempty"`
	MaximumAmount        *float64              `json:"maximum_amount,omitempty"`
	Amount               float64               `json:"amount"`
	CallbackURL          string                `json:"callback_url"`
	SenderTerminalCode   *string               `json:"sender_terminal_code,omitempty"`
	SenderTerminalData   interface{}           `json:"sender_terminal_data,omitempty"`
	AllowSubscribe       *bool                 `json:"allow_subscribe,omitempty"`
	SubscriptionInterval string                `json:"subscription_interval,omitempty"`
	SubscriptionWebhook  string                `json:"subscription_webhook,omitempty"`
	Note                 *string               `json:"note,omitempty"`
	Transactions         []Transaction         `json:"transactions,omitempty"`
	Lines                []InvoiceLine         `json:"lines,omitempty"`
}

// CreateSimpleInvoiceRequest is the request body for creating a simple invoice.
type CreateSimpleInvoiceRequest struct {
	InvoiceCode         string  `json:"invoice_code"`
	SenderInvoiceNo     string  `json:"sender_invoice_no"`
	InvoiceReceiverCode string  `json:"invoice_receiver_code"`
	InvoiceDescription  string  `json:"invoice_description"`
	SenderBranchCode    string  `json:"sender_branch_code,omitempty"`
	Amount              float64 `json:"amount"`
	CallbackURL         string  `json:"callback_url"`
}

// CreateEbarimtInvoiceRequest is the request body for creating an invoice with ebarimt.
type CreateEbarimtInvoiceRequest struct {
	InvoiceCode         string               `json:"invoice_code"`
	SenderInvoiceNo     string               `json:"sender_invoice_no"`
	SenderBranchCode    string               `json:"sender_branch_code,omitempty"`
	SenderStaffData     *SenderStaffData     `json:"sender_staff_data,omitempty"`
	SenderStaffCode     string               `json:"sender_staff_code,omitempty"`
	InvoiceReceiverCode string               `json:"invoice_receiver_code"`
	InvoiceReceiverData *InvoiceReceiverData `json:"invoice_receiver_data,omitempty"`
	InvoiceDescription  string               `json:"invoice_description"`
	TaxType             string               `json:"tax_type"`
	DistrictCode        string               `json:"district_code"`
	CallbackURL         string               `json:"callback_url"`
	Lines               []EbarimtInvoiceLine `json:"lines"`
}

// InvoiceResponse is the response from creating an invoice.
type InvoiceResponse struct {
	InvoiceID     string     `json:"invoice_id"`
	QRText        string     `json:"qr_text"`
	QRImage       string     `json:"qr_image"`
	QPay_ShortURL string     `json:"qPay_shortUrl"`
	URLs          []Deeplink `json:"urls"`
}

// --- Payment ---

// Offset represents pagination parameters.
type Offset struct {
	PageNumber int `json:"page_number"`
	PageLimit  int `json:"page_limit"`
}

// PaymentCheckRequest is the request body for checking a payment.
type PaymentCheckRequest struct {
	ObjectType string  `json:"object_type"`
	ObjectID   string  `json:"object_id"`
	Offset     *Offset `json:"offset,omitempty"`
}

// PaymentCheckResponse is the response from checking a payment.
type PaymentCheckResponse struct {
	Count      int               `json:"count"`
	PaidAmount float64           `json:"paid_amount,omitempty"`
	Rows       []PaymentCheckRow `json:"rows"`
}

// PaymentCheckRow represents a single payment check result row.
type PaymentCheckRow struct {
	PaymentID           string            `json:"payment_id"`
	PaymentStatus       string            `json:"payment_status"`
	PaymentAmount       string            `json:"payment_amount"`
	TrxFee              string            `json:"trx_fee"`
	PaymentCurrency     string            `json:"payment_currency"`
	PaymentWallet       string            `json:"payment_wallet"`
	PaymentType         string            `json:"payment_type"`
	NextPaymentDate     *string           `json:"next_payment_date"`
	NextPaymentDatetime *string           `json:"next_payment_datetime"`
	CardTransactions    []CardTransaction `json:"card_transactions"`
	P2PTransactions     []P2PTransaction  `json:"p2p_transactions"`
}

// PaymentDetail represents detailed payment information.
type PaymentDetail struct {
	PaymentID           string            `json:"payment_id"`
	PaymentStatus       string            `json:"payment_status"`
	PaymentFee          string            `json:"payment_fee"`
	PaymentAmount       string            `json:"payment_amount"`
	PaymentCurrency     string            `json:"payment_currency"`
	PaymentDate         string            `json:"payment_date"`
	PaymentWallet       string            `json:"payment_wallet"`
	TransactionType     string            `json:"transaction_type"`
	ObjectType          string            `json:"object_type"`
	ObjectID            string            `json:"object_id"`
	NextPaymentDate     *string           `json:"next_payment_date"`
	NextPaymentDatetime *string           `json:"next_payment_datetime"`
	CardTransactions    []CardTransaction `json:"card_transactions"`
	P2PTransactions     []P2PTransaction  `json:"p2p_transactions"`
}

// CardTransaction represents a card payment transaction.
type CardTransaction struct {
	CardMerchantCode     string `json:"card_merchant_code,omitempty"`
	CardTerminalCode     string `json:"card_terminal_code,omitempty"`
	CardNumber           string `json:"card_number,omitempty"`
	CardType             string `json:"card_type"`
	IsCrossBorder        bool   `json:"is_cross_border"`
	Amount               string `json:"amount,omitempty"`
	TransactionAmount    string `json:"transaction_amount,omitempty"`
	Currency             string `json:"currency,omitempty"`
	TransactionCurrency  string `json:"transaction_currency,omitempty"`
	Date                 string `json:"date,omitempty"`
	TransactionDate      string `json:"transaction_date,omitempty"`
	Status               string `json:"status,omitempty"`
	TransactionStatus    string `json:"transaction_status,omitempty"`
	SettlementStatus     string `json:"settlement_status"`
	SettlementStatusDate string `json:"settlement_status_date"`
}

// P2PTransaction represents a peer-to-peer payment transaction.
type P2PTransaction struct {
	TransactionBankCode string `json:"transaction_bank_code"`
	AccountBankCode     string `json:"account_bank_code"`
	AccountBankName     string `json:"account_bank_name"`
	AccountNumber       string `json:"account_number"`
	Status              string `json:"status"`
	Amount              string `json:"amount"`
	Currency            string `json:"currency"`
	SettlementStatus    string `json:"settlement_status"`
}

// PaymentListRequest is the request body for listing payments.
type PaymentListRequest struct {
	ObjectType string `json:"object_type"`
	ObjectID   string `json:"object_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Offset     Offset `json:"offset"`
}

// PaymentListResponse is the response from listing payments.
type PaymentListResponse struct {
	Count int               `json:"count"`
	Rows  []PaymentListItem `json:"rows"`
}

// PaymentListItem represents a single payment in a list response.
type PaymentListItem struct {
	PaymentID          string `json:"payment_id"`
	PaymentDate        string `json:"payment_date"`
	PaymentStatus      string `json:"payment_status"`
	PaymentFee         string `json:"payment_fee"`
	PaymentAmount      string `json:"payment_amount"`
	PaymentCurrency    string `json:"payment_currency"`
	PaymentWallet      string `json:"payment_wallet"`
	PaymentName        string `json:"payment_name"`
	PaymentDescription string `json:"payment_description"`
	QRCode             string `json:"qr_code"`
	PaidBy             string `json:"paid_by"`
	ObjectType         string `json:"object_type"`
	ObjectID           string `json:"object_id"`
}

// PaymentCancelRequest is the request body for canceling a payment.
type PaymentCancelRequest struct {
	CallbackURL string `json:"callback_url,omitempty"`
	Note        string `json:"note,omitempty"`
}

// PaymentRefundRequest is the request body for refunding a payment.
type PaymentRefundRequest struct {
	CallbackURL string `json:"callback_url,omitempty"`
	Note        string `json:"note,omitempty"`
}

// --- Ebarimt ---

// CreateEbarimtRequest is the request body for creating an ebarimt.
type CreateEbarimtRequest struct {
	PaymentID           string `json:"payment_id"`
	EbarimtReceiverType string `json:"ebarimt_receiver_type"`
	EbarimtReceiver     string `json:"ebarimt_receiver,omitempty"`
	DistrictCode        string `json:"district_code,omitempty"`
	ClassificationCode  string `json:"classification_code,omitempty"`
}

// EbarimtResponse is the response from creating or canceling an ebarimt.
type EbarimtResponse struct {
	ID                   string           `json:"id"`
	EbarimtBy            string           `json:"ebarimt_by"`
	GWalletID            string           `json:"g_wallet_id"`
	GWalletCustomerID    string           `json:"g_wallet_customer_id"`
	EbarimtReceiverType  string           `json:"ebarimt_receiver_type"`
	EbarimtReceiver      string           `json:"ebarimt_receiver"`
	EbarimtDistrictCode  string           `json:"ebarimt_district_code"`
	EbarimtBillType      string           `json:"ebarimt_bill_type"`
	GMerchantID          string           `json:"g_merchant_id"`
	MerchantBranchCode   string           `json:"merchant_branch_code"`
	MerchantTerminalCode *string          `json:"merchant_terminal_code"`
	MerchantStaffCode    *string          `json:"merchant_staff_code"`
	MerchantRegisterNo   string           `json:"merchant_register_no"`
	GPaymentID           string           `json:"g_payment_id"`
	PaidBy               string           `json:"paid_by"`
	ObjectType           string           `json:"object_type"`
	ObjectID             string           `json:"object_id"`
	Amount               string           `json:"amount"`
	VatAmount            string           `json:"vat_amount"`
	CityTaxAmount        string           `json:"city_tax_amount"`
	EbarimtQRData        string           `json:"ebarimt_qr_data"`
	EbarimtLottery       string           `json:"ebarimt_lottery"`
	Note                 *string          `json:"note"`
	BarimtStatus         string           `json:"barimt_status"`
	BarimtStatusDate     string           `json:"barimt_status_date"`
	EbarimtSentEmail     *string          `json:"ebarimt_sent_email"`
	EbarimtReceiverPhone string           `json:"ebarimt_receiver_phone"`
	TaxType              string           `json:"tax_type"`
	MerchantTIN          string           `json:"merchant_tin,omitempty"`
	EbarimtReceiptID     string           `json:"ebarimt_receipt_id,omitempty"`
	CreatedBy            string           `json:"created_by"`
	CreatedDate          string           `json:"created_date"`
	UpdatedBy            string           `json:"updated_by"`
	UpdatedDate          string           `json:"updated_date"`
	Status               bool             `json:"status"`
	BarimtItems          []EbarimtItem    `json:"barimt_items,omitempty"`
	BarimtTransactions   []interface{}    `json:"barimt_transactions,omitempty"`
	BarimtHistories      []EbarimtHistory `json:"barimt_histories,omitempty"`
}

// EbarimtItem represents a single item in an ebarimt receipt.
type EbarimtItem struct {
	ID                  string  `json:"id"`
	BarimtID            string  `json:"barimt_id"`
	MerchantProductCode *string `json:"merchant_product_code"`
	TaxProductCode      string  `json:"tax_product_code"`
	BarCode             *string `json:"bar_code"`
	Name                string  `json:"name"`
	UnitPrice           string  `json:"unit_price"`
	Quantity            string  `json:"quantity"`
	Amount              string  `json:"amount"`
	CityTaxAmount       string  `json:"city_tax_amount"`
	VatAmount           string  `json:"vat_amount"`
	Note                *string `json:"note"`
	CreatedBy           string  `json:"created_by"`
	CreatedDate         string  `json:"created_date"`
	UpdatedBy           string  `json:"updated_by"`
	UpdatedDate         string  `json:"updated_date"`
	Status              bool    `json:"status"`
}

// EbarimtHistory represents a historical record of an ebarimt.
type EbarimtHistory struct {
	ID                   string  `json:"id"`
	BarimtID             string  `json:"barimt_id"`
	EbarimtReceiverType  string  `json:"ebarimt_receiver_type"`
	EbarimtReceiver      string  `json:"ebarimt_receiver"`
	EbarimtRegisterNo    *string `json:"ebarimt_register_no"`
	EbarimtBillID        string  `json:"ebarimt_bill_id"`
	EbarimtDate          string  `json:"ebarimt_date"`
	EbarimtMacAddress    string  `json:"ebarimt_mac_address"`
	EbarimtInternalCode  string  `json:"ebarimt_internal_code"`
	EbarimtBillType      string  `json:"ebarimt_bill_type"`
	EbarimtQRData        string  `json:"ebarimt_qr_data"`
	EbarimtLottery       string  `json:"ebarimt_lottery"`
	EbarimtLotteryMsg    *string `json:"ebarimt_lottery_msg"`
	EbarimtErrorCode     *string `json:"ebarimt_error_code"`
	EbarimtErrorMsg      *string `json:"ebarimt_error_msg"`
	EbarimtResponseCode  *string `json:"ebarimt_response_code"`
	EbarimtResponseMsg   *string `json:"ebarimt_response_msg"`
	Note                 *string `json:"note"`
	BarimtStatus         string  `json:"barimt_status"`
	BarimtStatusDate     string  `json:"barimt_status_date"`
	EbarimtSentEmail     *string `json:"ebarimt_sent_email"`
	EbarimtReceiverPhone string  `json:"ebarimt_receiver_phone"`
	TaxType              string  `json:"tax_type"`
	CreatedBy            string  `json:"created_by"`
	CreatedDate          string  `json:"created_date"`
	UpdatedBy            string  `json:"updated_by"`
	UpdatedDate          string  `json:"updated_date"`
	Status               bool    `json:"status"`
}

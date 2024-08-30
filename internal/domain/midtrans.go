package domain

type MidtransCallback struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
}
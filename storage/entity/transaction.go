package entity

type Transaction struct {
	Id int64 `json:"id"`
	UserID int64 `json:"user_id"`
	AccountID int64 `json:"account_id"`
	CategoryID int64 `json:"category_id"`

	Type string `json:"type"`
	Amount float64 `json:"amount"`
	Payee string `json:"payee"`
	Currency string `json:"currency"`
	TransactionDate string `json:"transaction_date"`
	Description string `json:"description"`
	ReceiptURL string `json:"receipt_url"`
	Location string `json:"location"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
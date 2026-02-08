package entity

type Account struct {
	Id              int64   `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Currency        string  `json:"currency"`
	CurrentBalance  float64 `json:"current_balance"`
	BankName        string  `json:"bank_name"`
	AccountNumber   string  `json:"account_number"`
	IsIncludedInTotal bool    `json:"is_included_in_total"`
	UserID          int64   `json:"user_id"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
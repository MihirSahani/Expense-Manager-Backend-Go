package entity

type Category struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Type string `json:"type"`
	Color string `json:"color"`
}
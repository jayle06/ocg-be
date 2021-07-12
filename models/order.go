package models

type Order struct {
	ID          int64       `json:"id"`
	FullName    string      `json:"full_name"`
	PhoneNumber string      `json:"phone_number"`
	Email       string      `json:"email"`
	Address     string      `json:"address"`
	OrderItems  []OrderItem `json:"order_items"`
	Total       int64       `json:"total"`
	Payment     string      `json:"payment"`
	CreatedAt   string      `json:"created_at"`
}

package models

type Product struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       int64      `json:"price"`
	IsSale      bool       `json:"is_sale"`
	PriceSale   int64      `json:"price_sale"`
	Images      []Image    `json:"images"`
	Categories  []Category `json:"categories"`
	CreatedAt   string     `json:"created_at"`
}

package schema

type Product struct {
	ProductID       string  `json:"productID"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	Currency        string  `json:"currency"`
	Description     string  `json:"description"`
	ImageUrl        string  `json:"imageUrl"`
	CreateTimestamp string  `json:"createTimestamp"`
	UpdateTimestamp string  `json:"updateTimestamp"`
}

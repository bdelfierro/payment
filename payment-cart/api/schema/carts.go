package schema

const (
	CartStatusActive = "ACTIVE"
	CartStatusPAID   = "PAID"
)

// Carts represents schema of carts table
type Carts struct {
	UserID          string `json:"userID"`
	CartID          string `json:"cartID"`
	Status          string `json:"status"`
	CreateTimestamp string `json:"createTimestamp"`
	UpdateTimestamp string `json:"updateTimestamp"`
}

// CartItems represents schema of cart_items table
type CartItems struct {
	CartID          string `json:"cartID"`
	ProductID       string `json:"productID"`
	Quantity        int64  `json:"quantity"`
	Status          string `json:"status"`
	CreateTimestamp string `json:"createTimestamp"`
	UpdateTimestamp string `json:"updateTimestamp"`
}

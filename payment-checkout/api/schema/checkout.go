package schema

type LineItem struct {
	ProductName string  `json:"productName"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
	Currency    string  `json:"currency"`
	Quantity    int64   `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
}

type CheckoutRequest struct {
	CartID    string     `json:"cartID"`
	LineItems []LineItem `json:"lineItems"`
}

type CheckoutResponse struct {
	CheckoutRequest
	SessionID  string `json:"sessionID"`
	SessionURL string `json:"sessionUrl"`
}

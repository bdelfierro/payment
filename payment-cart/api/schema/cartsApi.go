package schema

type CartItem struct {
	ProductDetails Product `json:"productDetails"`
	Quantity       int64   `json:"quantity"`
}

type CartRequest struct {
	Action    string     `json:"action"`
	UserID    string     `json:"userID"`
	CartID    string     `json:"cartID"`
	LineItems []CartItem `json:"lineItems"`
}

type CartResponse struct {
	CartRequest
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

type CartStatusRequest struct {
	UserID string `json:"userID"`
	CartID string `json:"cartID"`
	Status string `json:"status"`
}

type CartStatusResponse struct {
	CartStatusRequest
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

type CartItemPrice struct {
	ProductDetails Product `json:"productDetails"`
	Quantity       int64   `json:"quantity"`
	TotalPrice     float64 `json:"totalPrice"`
}

type GetCartDetailsRequest struct {
	CartID string `json:"cartID"`
	UserID string `json:"userID"`
}

type GetCartDetailsResponse struct {
	GetCartDetailsRequest
	Status          string          `json:"status"`
	TotalAmount     float64         `json:"totalAmount"`
	TotalCount      int64           `json:"totalCount"`
	LineItems       []CartItemPrice `json:"lineItems"`
	ResponseCode    string          `json:"responseCode"`
	ResponseMessage string          `json:"responseMessage"`
}

package schema

// StripeSecret schema from google cloud secret manager for stripe api key
type StripeSecret struct {
	SecretKey string `json:"secretKey"`
}

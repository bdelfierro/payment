package schema

type StripeSecret struct {
	SecretKey        string `json:"secretKey"`
	WebhookSecretKey string `json:"webhookSecretKey"`
}

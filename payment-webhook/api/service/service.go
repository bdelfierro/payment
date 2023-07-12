package service

import (
	"context"
	"github.com/benedictdelfierro/payment/payment-webhook/api/messaging"
	"os"
)

type PaymentWebhook struct {
	MsgClient messaging.MessageClient
}

func NewPaymentWebhook(client messaging.MessageClient) *PaymentWebhook {
	return &PaymentWebhook{MsgClient: client}
}

func (s *PaymentWebhook) ProcessStripeEvent(ctx context.Context, sessionRawData []byte) error {

	msg := messaging.Message{
		Topic:          os.Getenv("FULFILLMENT_TOPIC_ID"),
		SubscriptionID: os.Getenv("FULFILLMENT_SUBS_ID"),
		Value:          sessionRawData,
	}

	s.MsgClient.Publish(ctx, msg)

	return nil
}

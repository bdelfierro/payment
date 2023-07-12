package messaging

import "context"

type Message struct {
	Topic          string
	SubscriptionID string
	Value          []byte
}

type MessageClient interface {
	Publish(ctx context.Context, message Message)
	//Subscribe(ctx context.Context, subscriptionID string)
}

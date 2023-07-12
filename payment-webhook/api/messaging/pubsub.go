package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
	"time"
)

type PubSubClient struct {
	C *pubsub.Client
}

func NewPubsubClient(ctx context.Context, projectID string) (*PubSubClient, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &PubSubClient{C: client}, nil
}

// Publish is the implementation of MessageClient interface for google PubSub
func (p *PubSubClient) Publish(ctx context.Context, message Message) {

	t := p.C.Topic(message.Topic)

	start := time.Now()
	result := t.Publish(ctx, &pubsub.Message{
		Data: message.Value,
	})
	log.Printf("async publish elapsed time : %v\n", time.Since(start))
	log.Printf("async publish result : %+v\n", result)

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		log.Printf("error getting result %v\n", err)
	}

	log.Printf("Published a message; msg ID: %v\n", id)

	log.Printf("sync publish elapsed time : %v\n", time.Since(start))

}

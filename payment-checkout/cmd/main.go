package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/benedictdelfierro/payment/payment-checkout/api/handlers"
	"github.com/benedictdelfierro/payment/payment-checkout/api/schema"
	"github.com/rs/cors"
	"github.com/stripe/stripe-go/v74"
	"hash/crc32"
	"log"
	"net/http"
	"os"
)

func main() {

	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	stripeSecretName := os.Getenv("STRIPE_SECRETS_NAME")
	responseURL := os.Getenv("STRIPE_RESPONSE_URL")
	secretsReq := fmt.Sprintf("projects/%v/secrets/%v/versions/latest", projectID, stripeSecretName)

	payload, err := accessSecretVersion(ctx, secretsReq)

	if err != nil {
		log.Fatalf("error accessing secrets manager, %v\n", err)
	}

	secretKey := schema.StripeSecret{}

	if err := json.Unmarshal(payload, &secretKey); err != nil {
		log.Fatalf("error parsing secrets payload, %v\n", err)
	}

	stripe.Key = secretKey.SecretKey

	handler := handlers.NewHandler(responseURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/create-checkout-session", handler.HandleCreateCheckoutSession)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Use the cors middleware with default options
	h := cors.Default().Handler(mux)

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), h); err != nil {
		log.Fatal(err)
	}

}

func accessSecretVersion(ctx context.Context, name string) ([]byte, error) {
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return nil, fmt.Errorf("Data corruption detected.")
	}

	// WARNING: Do not print the secret in a production environment - this snippet
	// is showing how to access the secret material.
	//log.Printf("Plaintext: %s\n", string(result.Payload.Data))

	return result.Payload.Data, nil
}

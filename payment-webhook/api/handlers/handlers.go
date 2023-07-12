package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/benedictdelfierro/payment/payment-webhook/api/messaging"
	"github.com/benedictdelfierro/payment/payment-webhook/api/service"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	StripeWebhookSecret string
	MsgClient           messaging.MessageClient
}

func NewHandler(webhookSecret string, msgClient messaging.MessageClient) *Handler {
	return &Handler{
		StripeWebhookSecret: webhookSecret,
		MsgClient:           msgClient,
	}
}

func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("io.ReadAll: %v", err)
		return
	}

	log.Printf("received msg from webhook: %+v\n", string(b))

	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// You can find your endpoint's secret in your webhook settings
	event, err := webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), h.StripeWebhookSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("webhook.ConstructEvent: %v", err)
		return
	}

	if event.Type == "checkout.session.completed" {
		var sessionData stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &sessionData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error parsing webhook JSON: %v", err)
			return
		}

		//log.Printf("sessionData: %+v\n", sessionData)

		//fmt.Println("Checkout Session completed, publishing message to pubsub")
		if err := service.NewPaymentWebhook(h.MsgClient).ProcessStripeEvent(ctx, b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error in processing stripe event: %v", err)
			return
		}
	} else {
		log.Printf("Unhandled event type: %v\n", event.Type)
	}

	writeJSON(w, "response from webhook")

}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}

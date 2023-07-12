package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/benedictdelfierro/payment/payment-checkout/api/schema"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"log"
	"net/http"
)

type Handler struct {
	RedirectUrl string
}

func NewHandler(RedirectUrl string) *Handler {
	return &Handler{
		RedirectUrl: RedirectUrl,
	}
}

// HandleCreateCheckoutSession handles request for session creation from client (e.g. flutter mobile app)
// It returns the stripe checkout session url
func (h *Handler) HandleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request from client")
	//ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	req := schema.CheckoutRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, item := range req.LineItems {
		item := &stripe.CheckoutSessionLineItemParams{
			Quantity: stripe.Int64(item.Quantity),
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:          stripe.String(item.Currency),
				UnitAmountDecimal: stripe.Float64(item.UnitPrice * 100),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Description: stripe.String(item.Description),
					Name:        stripe.String(item.ProductName),
					Images:      []*string{stripe.String(item.ImageURL)},
				},
			},
		}
		lineItems = append(lineItems, item)
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL:        stripe.String(h.RedirectUrl + "/success.html?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:         stripe.String(h.RedirectUrl + "/canceled.html"),
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:         lineItems,
		ClientReferenceID: stripe.String(req.CartID),
	}
	s, err := session.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("error while creating session %v", err.Error()), http.StatusInternalServerError)
		return
	}

	log.Printf("created session: %+v\n", s)
	resp := schema.CheckoutResponse{
		CheckoutRequest: schema.CheckoutRequest{
			CartID: req.CartID,
		},
		SessionID:  s.ID,
		SessionURL: s.URL,
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response %v", err.Error()), http.StatusInternalServerError)
		return
	}
	return

}

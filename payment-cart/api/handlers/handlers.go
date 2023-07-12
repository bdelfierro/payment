package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/benedictdelfierro/payment/payment-cart/api/schema"
	"github.com/benedictdelfierro/payment/payment-cart/api/service"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
)

type Handler struct {
	Service *service.CartService
}

func NewHandler(service *service.CartService) *Handler {
	return &Handler{
		Service: service,
	}
}

// HandleAddItem handles adding new item in the cart
func (h *Handler) HandleAddItem(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleAddItem")
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	req := schema.CartRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("error in parsing request, %v", err), http.StatusBadRequest)
		return
	}

	log.Printf(" req: %+v\n", req)

	response := schema.GetCartDetailsResponse{}

	resp, err := h.Service.AddItem(ctx, req)
	if err != nil {
		response.ResponseCode = "E001"
		response.ResponseMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	response.LineItems = resp.LineItems
	response.UserID = resp.UserID
	response.CartID = resp.CartID
	response.Status = resp.Status
	total := 0.0
	var totalCount int64
	for _, item := range resp.LineItems {
		total += item.TotalPrice
		totalCount += item.Quantity
	}
	response.TotalAmount = total
	response.TotalCount = totalCount
	response.ResponseCode = "0000"
	response.ResponseMessage = "success"

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error encoding json response", http.StatusBadRequest)
		return
	}
}

// HandleRemoveItem handles removing an item from the cart
func (h *Handler) HandleRemoveItem(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleRemoveItem")
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	req := schema.CartRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("error in parsing request, %v", err), http.StatusBadRequest)
		return
	}

	log.Printf(" req: %+v\n", req)

	response := schema.GetCartDetailsResponse{}

	resp, err := h.Service.RemoveItem(ctx, req)
	if err != nil {
		response.ResponseCode = "E001"
		response.ResponseMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	response.LineItems = resp.LineItems
	response.UserID = resp.UserID
	response.CartID = resp.CartID
	response.Status = resp.Status
	total := 0.0
	var totalCount int64
	for _, item := range resp.LineItems {
		total += item.TotalPrice
		totalCount += item.Quantity
	}
	response.TotalAmount = total
	response.TotalCount = totalCount
	response.ResponseCode = "0000"
	response.ResponseMessage = "success"

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error encoding json response", http.StatusBadRequest)
		return
	}
}

// HandleRemoveProduct handles deleting entire product from the cart
func (h *Handler) HandleRemoveProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleRemoveProduct")
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	req := schema.CartRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("error in parsing request, %v", err), http.StatusBadRequest)
		return
	}

	log.Printf(" req: %+v\n", req)

	response := schema.GetCartDetailsResponse{}

	resp, err := h.Service.RemoveProductFromCart(ctx, req)
	if err != nil {
		response.ResponseCode = "E001"
		response.ResponseMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	response.LineItems = resp.LineItems
	response.UserID = resp.UserID
	response.CartID = resp.CartID
	response.Status = resp.Status
	total := 0.0
	var totalCount int64
	for _, item := range resp.LineItems {
		total += item.TotalPrice
		totalCount += item.Quantity
	}
	response.TotalAmount = total
	response.TotalCount = totalCount
	response.ResponseCode = "0000"
	response.ResponseMessage = "success"

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error encoding json response", http.StatusBadRequest)
		return
	}
}

// HandleUpdateCartStatus handles update of cart payment status and triggered from stripe webhook
func (h *Handler) HandleUpdateCartStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleUpdateCartStatus")
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	req := schema.CartStatusRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("error in parsing request, %v", err), http.StatusBadRequest)
		return
	}

	log.Printf(" req: %+v\n", req)

	response := schema.CartStatusResponse{
		CartStatusRequest: req,
		ResponseCode:      "0000",
		ResponseMessage:   "Success",
	}

	if err := h.Service.UpdateCartStatus(ctx, req.CartID, req.Status); err != nil {
		response.ResponseCode = "E001"
		response.ResponseMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error encoding json response", http.StatusBadRequest)
		return
	}
}

// HandleGetCartDetails handles retrieving cart details
func (h *Handler) HandleGetCartDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleGetCartDetails")
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	req := schema.GetCartDetailsRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("error in parsing request, %v", err), http.StatusBadRequest)
		return
	}
	log.Printf(" req: %+v\n", req)

	response := schema.GetCartDetailsResponse{
		GetCartDetailsRequest: req,
	}

	if resp, err := h.Service.GetCartDetails(ctx, req.CartID); err != nil {
		response.ResponseCode = "E001"
		response.ResponseMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response = resp
		response.ResponseCode = "0000"
		response.ResponseMessage = "Success"
		total := decimal.NewFromFloat(0.0)
		totalCount := decimal.NewFromInt(0)
		for _, item := range resp.LineItems {
			total = total.Add(decimal.NewFromFloat(item.TotalPrice))
			totalCount = totalCount.Add(decimal.NewFromInt(item.Quantity))
		}
		response.TotalAmount, _ = total.Round(2).Float64()
		response.TotalCount = totalCount.IntPart()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error encoding json response", http.StatusBadRequest)
		return
	}
}

// HandleGetActiveCartDetails handles retrieving of active cart details
func (h *Handler) HandleGetActiveCartDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleGetActiveCartDetails")
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	req := schema.GetCartDetailsRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("error in parsing request, %v", err), http.StatusBadRequest)
		return
	}
	log.Printf(" req: %+v\n", req)

	response := schema.GetCartDetailsResponse{
		GetCartDetailsRequest: req,
	}

	if resp, err := h.Service.GetActiveCartDetails(ctx, req.UserID); err != nil {
		if len(resp.ResponseCode) == 0 {
			response.ResponseCode = "E001"
			response.ResponseMessage = err.Error()
		} else {
			response.ResponseCode = resp.ResponseCode
			response.ResponseMessage = resp.ResponseMessage
		}
		if response.ResponseCode == "E002" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		response = resp
		response.ResponseCode = "0000"
		response.ResponseMessage = "Success"
		total := decimal.NewFromFloat(0.0)
		totalCount := decimal.NewFromInt(0)
		for _, item := range resp.LineItems {
			total = total.Add(decimal.NewFromFloat(item.TotalPrice))
			totalCount = totalCount.Add(decimal.NewFromInt(item.Quantity))
		}
		response.TotalAmount, _ = total.Round(2).Float64()
		response.TotalCount = totalCount.IntPart()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error encoding json response", http.StatusBadRequest)
		return
	}
}

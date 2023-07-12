package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/benedictdelfierro/payment/products/api/schema"
	"github.com/benedictdelfierro/payment/products/api/service"
	"log"
	"net/http"
)

type Handler struct {
	Service *service.ProductsService
}

func NewHandler(service *service.ProductsService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) HandleProductList(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleProductList")
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	req := schema.ProductRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("error in parsing request, %v", err), http.StatusBadRequest)
		return
	}
	log.Printf(" req: %+v\n", req)

	response := schema.ProductResponse{
		ProductRequest:  req,
		ResponseCode:    "0000",
		ResponseMessage: "Success",
	}

	if resp, err := h.Service.GetProductDetails(ctx); err != nil {
		response.ResponseCode = "E001"
		response.ResponseMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Products = resp.Products
	}

	log.Printf("response: %+v", response)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error encoding json response", http.StatusBadRequest)
		return
	}
}

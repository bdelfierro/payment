package service

import (
	"context"
	"errors"
	"github.com/benedictdelfierro/payment/payment-cart/api/schema"
	"github.com/benedictdelfierro/payment/payment-cart/storage"
)

type CartService struct {
	db storage.Carts
}

func NewCartService(db storage.Carts) *CartService {
	return &CartService{
		db,
	}
}

func (c *CartService) AddItem(ctx context.Context, req schema.CartRequest) (schema.GetCartDetailsResponse, error) {

	if len(req.LineItems) == 0 {
		return schema.GetCartDetailsResponse{}, errors.New("no item in the list")
	}

	item := schema.CartItem{
		ProductDetails: req.LineItems[0].ProductDetails,
		Quantity:       req.LineItems[0].Quantity,
	}

	return c.db.AddItem(ctx, req.UserID, item)

}

func (c *CartService) RemoveItem(ctx context.Context, req schema.CartRequest) (schema.GetCartDetailsResponse, error) {

	if len(req.LineItems) == 0 {
		return schema.GetCartDetailsResponse{}, errors.New("no item in the list")
	}

	item := schema.CartItem{
		ProductDetails: req.LineItems[0].ProductDetails,
		Quantity:       req.LineItems[0].Quantity,
	}

	return c.db.RemoveItem(ctx, req.UserID, req.CartID, item)

}

func (c *CartService) RemoveProductFromCart(ctx context.Context, req schema.CartRequest) (schema.GetCartDetailsResponse, error) {

	if len(req.LineItems) == 0 {
		return schema.GetCartDetailsResponse{}, errors.New("no item in the list")
	}

	item := schema.CartItem{
		ProductDetails: req.LineItems[0].ProductDetails,
		Quantity:       req.LineItems[0].Quantity,
	}

	return c.db.RemoveProductFromCart(ctx, req.CartID, item)

}

func (c *CartService) UpdateCartStatus(ctx context.Context, cartID string, status string) error {

	if len(status) == 0 {
		status = schema.CartStatusPAID
	}
	return c.db.UpdateCartStatus(ctx, cartID, status)

}

func (c *CartService) GetCartDetails(ctx context.Context, cartID string) (schema.GetCartDetailsResponse, error) {

	return c.db.GetCartDetails(ctx, cartID)

}

func (c *CartService) GetActiveCartDetails(ctx context.Context, userID string) (schema.GetCartDetailsResponse, error) {

	return c.db.GetActiveCartDetails(ctx, userID)

}

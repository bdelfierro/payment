package storage

import (
	"context"
	"github.com/benedictdelfierro/payment/payment-cart/api/schema"
)

// Carts interface to handle query and update on shopping cart
type Carts interface {
	AddItem(ctx context.Context, userID string, item schema.CartItem) (schema.GetCartDetailsResponse, error)
	RemoveItem(ctx context.Context, userID string, cartID string, item schema.CartItem) (schema.GetCartDetailsResponse, error)
	RemoveProductFromCart(ctx context.Context, cartID string, item schema.CartItem) (schema.GetCartDetailsResponse, error)
	UpdateCartStatus(ctx context.Context, cartID string, status string) error
	GetCartDetails(ctx context.Context, cartID string) (schema.GetCartDetailsResponse, error)
	GetActiveCartDetails(ctx context.Context, userID string) (schema.GetCartDetailsResponse, error)
}

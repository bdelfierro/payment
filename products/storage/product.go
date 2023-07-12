package storage

import (
	"context"
	"github.com/benedictdelfierro/payment/products/api/schema"
)

type Products interface {
	GetProductList(ctx context.Context) (schema.ProductResponse, error)
}

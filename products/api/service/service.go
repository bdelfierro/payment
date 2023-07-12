package service

import (
	"context"
	"github.com/benedictdelfierro/payment/products/api/schema"
	"github.com/benedictdelfierro/payment/products/storage"
)

type ProductsService struct {
	db storage.Products
}

func NewProductService(db storage.Products) *ProductsService {
	return &ProductsService{
		db: db,
	}
}

func (p *ProductsService) GetProductDetails(ctx context.Context) (schema.ProductResponse, error) {

	return p.db.GetProductList(ctx)

}

package postgreSQL

import (
	"context"
	"database/sql"
	"github.com/benedictdelfierro/payment/products/api/schema"
	"log"
)

type ProductsSQL struct {
	DBConn *sql.DB
}

func NewCProductsPSQL(dbConn *sql.DB) *ProductsSQL {
	return &ProductsSQL{
		DBConn: dbConn,
	}
}

func (p *ProductsSQL) GetProductList(ctx context.Context) (schema.ProductResponse, error) {

	var resp schema.ProductResponse

	rows, err := p.DBConn.QueryContext(ctx, GetProductList)

	if err != nil {
		log.Printf("error selecting products table: %v\n", err)
		return resp, err
	}

	var products []schema.Product

	defer rows.Close()

	for rows.Next() {

		product := schema.Product{}

		if err := rows.Scan(&product.ProductID, &product.Name, &product.Price,
			&product.Currency, &product.Description, &product.ImageUrl, &product.CreateTimestamp, &product.UpdateTimestamp); err != nil {
			return resp, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return resp, err
	}
	resp.Products = products

	return resp, nil
}

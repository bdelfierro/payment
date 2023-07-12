package postgreSQL

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/benedictdelfierro/payment/payment-cart/api/schema"
	"github.com/shopspring/decimal"
	"log"
)

type CartsPSQL struct {
	DBConn *sql.DB
}

func NewCartsPSQL(dbConn *sql.DB) *CartsPSQL {
	return &CartsPSQL{
		DBConn: dbConn,
	}
}

func (c *CartsPSQL) AddItem(ctx context.Context, userID string, item schema.CartItem) (schema.GetCartDetailsResponse, error) {
	var cartID string
	resp := schema.GetCartDetailsResponse{}

	tx, err := c.DBConn.BeginTx(ctx, nil)
	if err != nil {
		return resp, err
	}
	err = tx.QueryRowContext(ctx, CheckActiveCart, userID).Scan(&cartID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("error checking active cart, %v\n", err)
			return resp, err
		}
		result, err := tx.ExecContext(ctx, CreateCart, userID)
		if err != nil {
			log.Printf("error creating new cart id, %v\n", err)
			if rollErr := tx.Rollback(); err != nil {
				return resp, fmt.Errorf("%v, rollback err: %v", err, rollErr)
			}
			return resp, err
		}
		rowCount, err := result.RowsAffected() // get the number of rows affected

		if err != nil {
			log.Fatal(err) // if there's an error retrieving the row count, log it and exit
		}
		fmt.Printf("inserted %d rows in carts table.\n", rowCount)

		err = tx.QueryRowContext(ctx, CheckActiveCart, userID).Scan(&cartID)
		if err != nil {
			log.Printf("error getting newly created cart id, %v\n", err)
			if rollErr := tx.Rollback(); err != nil {
				return resp, fmt.Errorf("%v, rollback err: %v", err, rollErr)
			}
			return resp, err
		}
	}

	result, err := tx.ExecContext(ctx, Upsert, cartID, item.ProductDetails.ProductID, item.Quantity)
	if err != nil {
		log.Printf("error creating new cart item, %v\n", err)
		if rollErr := tx.Rollback(); err != nil {
			return resp, fmt.Errorf("%v, rollback err: %v", err, rollErr)
		}
		return resp, err
	}

	rowCount, err := result.RowsAffected() // get the number of rows affected
	if err != nil {
		if rollErr := tx.Rollback(); err != nil {
			return resp, fmt.Errorf("%v, rollback err: %v", err, rollErr)
		}
		return resp, err
	}
	fmt.Printf("inserted %d rows in cart_items table.\n", rowCount)

	if err := tx.Commit(); err != nil {
		log.Printf("error in commit %v", err)
		return resp, err
	}

	if result, err := c.GetCartDetails(ctx, cartID); err != nil {
		return resp, err
	} else {
		resp = result
	}

	return resp, nil
}

func (c *CartsPSQL) RemoveItem(ctx context.Context, userID string, cartID string, item schema.CartItem) (schema.GetCartDetailsResponse, error) {

	resp := schema.GetCartDetailsResponse{}

	tx, err := c.DBConn.BeginTx(ctx, nil)
	if err != nil {
		return resp, err
	}

	result, err := tx.ExecContext(ctx, Upsert, cartID, item.ProductDetails.ProductID, item.Quantity)
	if err != nil {
		log.Printf("error creating new cart item, %v\n", err)
		if rollErr := tx.Rollback(); err != nil {
			return resp, fmt.Errorf("%v, rollback err: %v", err, rollErr)
		}
		return resp, err
	}

	rowCount, err := result.RowsAffected() // get the number of rows affected
	if err != nil {
		if rollErr := tx.Rollback(); err != nil {
			return resp, fmt.Errorf("%v, rollback err: %v", err, rollErr)
		}
		return resp, err
	}
	fmt.Printf("affected %d rows in cart_items table.\n", rowCount)

	result, err = tx.ExecContext(ctx, DeleteEmptyItems, cartID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("error deleting empty items, %v\n", err)
			return resp, err
		}
	}

	rowCount, err = result.RowsAffected() // get the number of rows affected
	if err != nil {
		if rollErr := tx.Rollback(); err != nil {
			return resp, fmt.Errorf("%v, rollback err: %v", err, rollErr)
		}
		return resp, err
	}
	fmt.Printf("affected deleted %d rows in cart_items table.\n", rowCount)

	if err := tx.Commit(); err != nil {
		log.Printf("error in commit %v", err)
		return resp, err
	}

	if result, err := c.GetCartDetails(ctx, cartID); err != nil {
		return resp, err
	} else {
		resp = result
	}

	return resp, nil
}

func (c *CartsPSQL) RemoveProductFromCart(ctx context.Context, cartID string, item schema.CartItem) (schema.GetCartDetailsResponse, error) {

	resp := schema.GetCartDetailsResponse{}

	tx, err := c.DBConn.BeginTx(ctx, nil)
	if err != nil {
		return resp, err
	}

	result, err := tx.ExecContext(ctx, DeleteProductFromCart, cartID, item.ProductDetails.ProductID)
	if err != nil {
		log.Printf("error deleting product %v from cart, %v\n", item.ProductDetails.ProductID, err)
		if rollErr := tx.Rollback(); err != nil {
			return resp, fmt.Errorf("%v, rollback err: %v", err, rollErr)
		}
		return resp, err
	}

	rowCount, err := result.RowsAffected() // get the number of rows affected
	if err != nil {
		if rollErr := tx.Rollback(); err != nil {
			return resp, fmt.Errorf("%v, rollback err: %v", err, rollErr)
		}
		return resp, err
	}
	fmt.Printf("affected %d rows in cart_items table.\n", rowCount)

	if err := tx.Commit(); err != nil {
		log.Printf("error in commit %v", err)
		return resp, err
	}

	if result, err := c.GetCartDetails(ctx, cartID); err != nil {
		return resp, err
	} else {
		resp = result
	}

	return resp, nil
}

func (c *CartsPSQL) UpdateCartStatus(ctx context.Context, cartID string, status string) error {

	tx, err := c.DBConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, UpdateCartStatus, status, cartID)

	if err != nil {
		log.Printf("error updating cart status, %v\n", err)
		if rollErr := tx.Rollback(); err != nil {
			return fmt.Errorf("%v, rollback err: %v", err, rollErr)
		}
		return err
	}

	rowCount, err := result.RowsAffected() // get the number of rows affected
	if err != nil {
		if rollErr := tx.Rollback(); err != nil {
			return fmt.Errorf("%v, rollback err: %v", err, rollErr)
		}
		return err
	}
	fmt.Printf("updated %d rows in carts table.\n", rowCount)

	if err := tx.Commit(); err != nil {
		log.Printf("error in commit %v", err)
		return err
	}
	log.Printf("updated cartID: %v to status: %v\n", cartID, status)

	return nil
}

func (c *CartsPSQL) GetCartDetails(ctx context.Context, cartID string) (schema.GetCartDetailsResponse, error) {

	resp := schema.GetCartDetailsResponse{}

	rows, err := c.DBConn.QueryContext(ctx, GetCartDetails, cartID)

	if err != nil {
		log.Printf("error selecting cartdetails for cardID: %v, %v\n", cartID, err)
		return resp, err
	}

	var lineItems []schema.CartItemPrice

	defer rows.Close()

	for rows.Next() {

		cartItem := schema.CartItemPrice{}

		if err := rows.Scan(&resp.UserID, &resp.CartID, &resp.Status, &cartItem.ProductDetails.ProductID,
			&cartItem.Quantity, &cartItem.ProductDetails.Name, &cartItem.ProductDetails.Price, &cartItem.ProductDetails.Currency, &cartItem.ProductDetails.Description,
			&cartItem.ProductDetails.ImageUrl, &cartItem.TotalPrice); err != nil {
			return resp, err
		}
		value := decimal.NewFromInt(cartItem.Quantity).Mul(decimal.NewFromFloat(cartItem.ProductDetails.Price))
		cartItem.TotalPrice, _ = value.Float64()
		lineItems = append(lineItems, cartItem)
	}

	if err := rows.Err(); err != nil {
		return resp, err
	}
	resp.LineItems = lineItems

	return resp, nil
}

func (c *CartsPSQL) GetActiveCartDetails(ctx context.Context, userID string) (schema.GetCartDetailsResponse, error) {
	resp := schema.GetCartDetailsResponse{}
	var cartID string

	if err := c.DBConn.QueryRowContext(ctx, CheckActiveCart, userID).Scan(&cartID); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("error checking active cart, %v\n", err)
			return resp, err
		} else {
			resp.ResponseCode = "E002"
			resp.ResponseMessage = "NO ACTIVE CART FOUND"
			return resp, err
		}
	}

	r, err := c.GetCartDetails(ctx, cartID)
	if err != nil {
		return resp, err
	} else {
		resp = r
	}

	return resp, nil
}

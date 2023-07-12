package main

import (
	"context"
	"fmt"
	"github.com/benedictdelfierro/payment/payment-cart/api/handlers"
	"github.com/benedictdelfierro/payment/payment-cart/api/service"
	"github.com/benedictdelfierro/payment/payment-cart/storage/postgreSQL"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {

	ctx := context.Background()

	// Initialize the database connection
	dbConn, err := postgreSQL.NewConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := dbConn.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	cartSQL := postgreSQL.NewCartsPSQL(dbConn)

	s := service.NewCartService(cartSQL)

	handler := handlers.NewHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("/addItem", handler.HandleAddItem)
	mux.HandleFunc("/removeItem", handler.HandleRemoveItem)
	mux.HandleFunc("/removeProduct", handler.HandleRemoveProduct)
	mux.HandleFunc("/updateCartStatus", handler.HandleUpdateCartStatus)
	mux.HandleFunc("/getCartDetails", handler.HandleGetCartDetails)
	mux.HandleFunc("/getActiveCartDetails", handler.HandleGetActiveCartDetails)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Use the cors middleware with default options
	h := cors.Default().Handler(mux)

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), h); err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"context"
	"fmt"
	"github.com/benedictdelfierro/payment/products/api/handlers"
	"github.com/benedictdelfierro/payment/products/api/service"
	"github.com/benedictdelfierro/payment/products/storage/postgreSQL"
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

	productsSQL := postgreSQL.NewCProductsPSQL(dbConn)

	s := service.NewProductService(productsSQL)

	handler := handlers.NewHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("/productlist", handler.HandleProductList)

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

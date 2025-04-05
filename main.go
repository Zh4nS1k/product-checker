package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"product-checker/handlers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/check-product", handlers.CheckProduct).Methods("POST")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

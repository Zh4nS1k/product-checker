package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"product-checker/database"
	"product-checker/handlers"
)

func main() {
	database.Connect()

	r := mux.NewRouter()
	r.HandleFunc("/api/check-product", handlers.CheckProduct).Methods("POST")
	r.HandleFunc("/api/history", handlers.GetHistory).Methods("GET")
	r.HandleFunc("/api/history/{id}", handlers.GetHistoryByID).Methods("GET")
	r.HandleFunc("/api/history/{id}", handlers.UpdateHistory).Methods("PUT")
	r.HandleFunc("/api/history/{id}", handlers.DeleteHistory).Methods("DELETE")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend/"))))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

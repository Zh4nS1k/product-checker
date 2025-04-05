package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"product-checker/database"
	"product-checker/models"
	"product-checker/utils"
	"time"
)

// CheckProduct: Handles the barcode check and creates a new history entry
func CheckProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Barcode string `json:"barcode"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	original := utils.IsBarcodeValid(input.Barcode)
	country := utils.GetCountryFromBarcode(input.Barcode)
	now := time.Now().Format(time.RFC3339)

	entry := models.CheckedProduct{
		Barcode:    input.Barcode,
		IsOriginal: original,
		Country:    country,
		CheckedAt:  now,
	}

	collection := database.GetCollection()
	result, _ := collection.InsertOne(r.Context(), entry)
	entry.ID = result.InsertedID.(primitive.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

// GetHistory: Fetch all product check history records
func GetHistory(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection()

	// Use bson.D{} to create an empty filter for fetching all documents
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var products []models.CheckedProduct
	if err := cursor.All(context.Background(), &products); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetHistoryByID: Fetch a product check history by ID
func GetHistoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection()
	var product models.CheckedProduct
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateHistory: Update a specific product check history by ID
func UpdateHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedProduct models.CheckedProduct
	json.NewDecoder(r.Body).Decode(&updatedProduct)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection()
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"barcode":     updatedProduct.Barcode,
			"is_original": updatedProduct.IsOriginal,
			"country":     updatedProduct.Country,
			"checked_at":  updatedProduct.CheckedAt,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

// DeleteHistory: Delete a product check history by ID
func DeleteHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection()
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

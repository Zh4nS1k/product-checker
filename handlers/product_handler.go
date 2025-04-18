package handlers

import (
	"context"
	"net/http"
	"product-checker/database"
	"product-checker/models"
	"product-checker/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckProduct(c *gin.Context) {
	var input struct {
		Barcode string `json:"barcode"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

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
	result, _ := collection.InsertOne(c.Request.Context(), entry)
	entry.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, entry)
}

func GetHistory(c *gin.Context) {
	collection := database.GetCollection()
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var products []models.CheckedProduct
	if err := cursor.All(context.Background(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetHistoryByID(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := database.GetCollection()
	var product models.CheckedProduct
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateHistory(c *gin.Context) {
	id := c.Param("id")

	var updatedProduct models.CheckedProduct
	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := database.GetCollection()
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"barcode":     updatedProduct.Barcode,
		"is_original": updatedProduct.IsOriginal,
		"country":     updatedProduct.Country,
		"checked_at":  updatedProduct.CheckedAt,
	}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func DeleteHistory(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := database.GetCollection()
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

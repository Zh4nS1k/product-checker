package handlers

import (
	"log"
	"net/http"
	"product-checker/database"
	"product-checker/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckProduct(c *gin.Context) {
	var input struct {
		ProductID string `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	isValid := utils.IsBarcodeValid(input.ProductID)
	result := "Valid"
	if !isValid {
		result = "Invalid"
	}

	// Получаем username из контекста (гарантировано middleware)
	username, _ := c.Get("username")
	log.Printf("Saving history for user: %s, product: %s", username.(string), input.ProductID)

	err := database.AddHistoryToPostgres(username.(string), input.ProductID, result)
	if err != nil {
		log.Printf("Failed to save history: %v", err)
	}

	country := utils.GetCountryFromBarcode(input.ProductID)

	c.JSON(http.StatusOK, gin.H{
		"product_id":    input.ProductID,
		"result":        result,
		"is_original":   isValid,
		"country":       country,
		"checked_at":    time.Now().Format(time.RFC3339),
		"history_saved": true,
	})
}

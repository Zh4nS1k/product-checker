package handlers

import (
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

	// Валидация штрих-кода
	isValid := utils.IsBarcodeValid(input.ProductID)
	result := "Valid"
	if !isValid {
		result = "Invalid"
	}

	country := utils.GetCountryFromBarcode(input.ProductID)

	responseData := gin.H{
		"product_id":  input.ProductID,
		"result":      result,
		"is_original": isValid,
		"country":     country,
		"checked_at":  time.Now().Format(time.RFC3339),
	}

	// Сохранение истории только для авторизованных пользователей
	if username, exists := c.Get("username"); exists {
		err := database.AddHistoryToMongo(username.(string), input.ProductID, result)
		if err != nil {
			responseData["history_error"] = "Failed to save history: " + err.Error()
		} else {
			responseData["history_saved"] = true
		}
	}

	c.JSON(http.StatusOK, responseData)
}

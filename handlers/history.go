package handlers

import (
	"net/http"
	"product-checker/database"
	"product-checker/models"

	"github.com/gin-gonic/gin"
)

func GetHistory(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
		return
	}

	history, err := database.GetHistoryByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get history",
			"details": err.Error(),
		})
		return
	}

	// Гарантируем, что возвращаем массив, даже если он пустой
	if history == nil {
		history = []models.ProductCheckHistory{}
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history, // Явно возвращаем поле "history" с массивом
	})
}

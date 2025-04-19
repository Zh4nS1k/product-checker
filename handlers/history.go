package handlers

import (
	"log"
	"net/http"
	"product-checker/database"
	"product-checker/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GetHistory возвращает историю проверок для авторизованного пользователя
// @Summary Получить историю проверок
// @Description Возвращает список всех проверенных товаров для текущего пользователя
// @Tags history
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} gin.H "История проверок"
// @Failure 401 {object} gin.H "Ошибка авторизации"
// @Failure 500 {object} gin.H "Ошибка сервера"
// @Router /api/history [get]
func GetHistory(c *gin.Context) {
	startTime := time.Now()
	username, exists := c.Get("username")
	if !exists {
		log.Printf("[%s] Unauthorized: username not in context", c.Request.Method)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authorization required",
			"success": false,
		})
		return
	}

	log.Printf("[%s] Fetching history for: %s", c.Request.Method, username.(string))

	history, err := database.GetHistoryByUsername(username.(string))
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve history",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	if history == nil {
		history = []models.CheckedProduct{} // Гарантируем возврат массива
	}

	log.Printf("[%s] History retrieved for %s in %v",
		c.Request.Method,
		username.(string),
		time.Since(startTime))

	c.JSON(http.StatusOK, gin.H{
		"count":    len(history),
		"history":  history,
		"success":  true,
		"duration": time.Since(startTime).String(),
	})
}

package main

import (
	"log"
	_ "net/http"
	"product-checker/database"
	"product-checker/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	// API маршруты
	api := r.Group("/api")
	{
		api.POST("/check-product", handlers.CheckProduct)
		api.GET("/history", handlers.GetHistory)
		api.GET("/history/:id", handlers.GetHistoryByID)
		api.PUT("/history/:id", handlers.UpdateHistory)
		api.DELETE("/history/:id", handlers.DeleteHistory)
	}

	// Отдаём статические файлы из папки frontend по пути /static
	r.Static("/static", "./frontend")

	// Главная страница — login.html
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	log.Println("🚀 Server started at :8080")
	r.Run(":8080")
}

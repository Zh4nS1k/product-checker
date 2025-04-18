package main

import (
	"log"
	"product-checker/database"
	"product-checker/handlers"
	"product-checker/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Отдаем статичные файлы из папки public
	r.Static("/static", "./public")

	// Статические страницы
	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})
	r.GET("/login", func(c *gin.Context) {
		c.File("./public/login.html")
	})
	r.GET("/register", func(c *gin.Context) {
		c.File("./public/register.html")
	})

	// Соединение с базами данных
	database.ConnectMongo()
	database.ConnectPostgres()

	// API маршруты
	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.POST("/check-product", handlers.CheckProduct)

		// Защищенные маршруты с авторизацией
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			protected.GET("/history", handlers.GetHistory)
			protected.GET("/history/:id", handlers.GetHistoryByID)
			protected.PUT("/history/:id", handlers.UpdateHistory)
			protected.DELETE("/history/:id", handlers.DeleteHistory)
		}
	}

	// Логирование старта сервера
	log.Println("🚀 Server started at :8080")
	r.Run(":8080")
}

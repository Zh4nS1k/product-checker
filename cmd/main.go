package main

import (
	"log"
	"product-checker/database"
	"product-checker/handlers"
	"product-checker/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	database.ConnectPostgres()
	database.ConnectMongo()

	r := gin.Default()

	// Статические файлы
	r.Static("/static", "./public")
	r.GET("/", func(c *gin.Context) { c.File("./public/index.html") })
	r.GET("/login", func(c *gin.Context) { c.File("./public/login.html") })
	r.GET("/register", func(c *gin.Context) { c.File("./public/register.html") })

	// API endpoints
	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.POST("/check-product", handlers.CheckProduct)
	}

	// Protected endpoints
	protected := api.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/history", handlers.GetHistory)
	}

	log.Println("🚀 Server started at :8080")
	r.Run(":8080")
}

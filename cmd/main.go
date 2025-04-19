package main

import (
	"log"
	"product-checker/database"
	"product-checker/handlers"
	"product-checker/middleware"
	"product-checker/migrations"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	database.ConnectPostgres()
	migrations.RunMigrations(database.DB)

	// Инициализация роутера
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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

		// Защищенные endpoints
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			protected.POST("/check-product", handlers.CheckProduct)
			protected.GET("/history", handlers.GetHistory)
		}
	}

	log.Println("🚀 Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

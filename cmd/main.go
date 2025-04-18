package main

import (
	"log"
	"product-checker/database"
	"product-checker/handlers"
	"product-checker/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.ConnectPostgres()

	r := gin.Default()

	// Static files
	r.Static("/static", "./public")

	// Pages
	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})
	r.GET("/login", func(c *gin.Context) {
		c.File("./public/login.html")
	})
	r.GET("/register", func(c *gin.Context) {
		c.File("./public/register.html")
	})

	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			protected.POST("/check-product", handlers.CheckProduct)
			protected.GET("/history", handlers.GetHistory)
			protected.GET("/history/:id", handlers.GetHistoryByID)
			protected.PUT("/history/:id", handlers.UpdateHistory)
			protected.DELETE("/history/:id", handlers.DeleteHistory)
		}
	}

	log.Println("🚀 Server started at :8080")
	r.Run(":8080")
}

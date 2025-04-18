package database

import (
	"log"
	"product-checker/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	dsn := "host=localhost user=postgres password=123456 dbname=product_checker port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}

	log.Println("✅ Connected to PostgreSQL")
}

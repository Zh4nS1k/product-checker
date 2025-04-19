package database

import (
	"context"
	"log"
	"product-checker/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectPostgres() {
	dsn := "host=localhost user=postgres password=123456 dbname=product_checker port=5432 sslmode=disable"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get SQL DB:", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ PostgreSQL connection established")
}

func AddHistoryToPostgres(username, productID, result string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	record := models.CheckedProduct{
		Username:  username,
		ProductID: productID,
		Result:    result,
		CheckedAt: time.Now(),
	}

	if err := DB.WithContext(ctx).Create(&record).Error; err != nil {
		log.Printf("Failed to save history: %v", err)
		return err
	}

	log.Printf("History saved - ID: %d, User: %s", record.ID, username)
	return nil
}

func GetHistoryByUsername(username string) ([]models.CheckedProduct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var history []models.CheckedProduct

	err := DB.WithContext(ctx).
		Where("username = ?", username).
		Order("checked_at DESC").
		Find(&history).Error

	if err != nil {
		return nil, err
	}

	return history, nil
}

package migrations

import (
	"log"
	"product-checker/models"

	"gorm.io/gorm"
)

// RunMigrations выполняет все необходимые миграции базы данных
func RunMigrations(db *gorm.DB) {
	tables := []interface{}{
		&models.User{},
		&models.CheckedProduct{},
	}

	// Удаляем и создаем таблицы заново (только для разработки!)
	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("⚠️ Failed to drop table: %v", err)
		}
	}

	// Автомиграции
	if err := db.AutoMigrate(tables...); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// Создаем индексы
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_checked_products_username 
		ON checked_products(username)
	`).Error; err != nil {
		log.Printf("⚠️ Failed to create username index: %v", err)
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_checked_products_checked_at 
		ON checked_products(checked_at DESC)
	`).Error; err != nil {
		log.Printf("⚠️ Failed to create checked_at index: %v", err)
	}
	// Добавить в функцию RunMigrations
	if err := db.AutoMigrate(&models.CheckedProduct{}); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// Добавляем столбец notes если его нет
	if !db.Migrator().HasColumn(&models.CheckedProduct{}, "notes") {
		if err := db.Migrator().AddColumn(&models.CheckedProduct{}, "notes"); err != nil {
			log.Printf("⚠️ Failed to add notes column: %v", err)
		}
	}

	log.Println("✅ Database migrations completed successfully")
}

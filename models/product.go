package models

import (
	"time"

	"gorm.io/gorm"
)

// CheckedProduct представляет запись о проверке товара
type CheckedProduct struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"index;size:100;not null" json:"username"`
	ProductID string         `gorm:"size:50;not null" json:"product_id"`
	Result    string         `gorm:"size:20;not null" json:"result"`
	Notes     string         `gorm:"size:500" json:"notes"` // Новое поле для заметок
	CheckedAt time.Time      `gorm:"index;not null" json:"checked_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName задает имя таблицы в БД
func (CheckedProduct) TableName() string {
	return "checked_products"
}

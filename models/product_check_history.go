package models

import (
	"time"
)

type ProductCheckHistory struct {
	Username  string    `bson:"username"`
	ProductID string    `bson:"product_id"`
	Result    string    `bson:"result"`
	CheckedAt time.Time `bson:"checked_at"`
}

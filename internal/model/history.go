package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ProductCheck struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserID         uint               `bson:"user_id"`
	Barcode        string             `bson:"barcode"`
	IsOriginal     bool               `bson:"is_original"`
	BarcodeCountry string             `bson:"barcode_country"`
	IPCountry      string             `bson:"ip_country"`
	CheckedAt      time.Time          `bson:"checked_at"`
}

type ProductCheckResult struct {
	ID             string    `json:"id,omitempty"`
	Barcode        string    `json:"barcode"`
	IsOriginal     bool      `json:"is_original"`
	BarcodeCountry string    `json:"barcode_country"`
	IPCountry      string    `json:"ip_country,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}

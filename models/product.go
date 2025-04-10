package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CheckedProduct struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Barcode    string             `json:"barcode"`
	IsOriginal bool               `json:"is_original"`
	Country    string             `json:"country"`
	CheckedAt  string             `json:"checked_at"`
}

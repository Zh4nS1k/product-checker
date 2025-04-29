package repository

import (
	"barcode-checker/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HistoryRepository interface {
	Create(check *model.ProductCheck) error
	GetByUserID(userID uint, page, limit int) ([]model.ProductCheck, int64, error)
	DeleteByID(userID uint, id string) error
}

type historyRepository struct {
	collection *mongo.Collection
}

func NewHistoryRepository(db *mongo.Database) HistoryRepository {
	return &historyRepository{
		collection: db.Collection("product_checks"),
	}
}

func (r *historyRepository) Create(check *model.ProductCheck) error {
	_, err := r.collection.InsertOne(context.Background(), check)
	return err
}

func (r *historyRepository) GetByUserID(userID uint, page, limit int) ([]model.ProductCheck, int64, error) {
	skip := (page - 1) * limit
	filter := bson.M{"user_id": userID}

	total, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"checked_at", -1}})

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var checks []model.ProductCheck
	if err = cursor.All(context.Background(), &checks); err != nil {
		return nil, 0, err
	}

	return checks, total, nil
}

func (r *historyRepository) DeleteByID(userID uint, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "user_id": userID}
	_, err = r.collection.DeleteOne(context.Background(), filter)
	return err
}

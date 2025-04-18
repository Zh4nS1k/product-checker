package database

import (
	"context"
	"log"
	"product-checker/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
)

func ConnectMongo() {
	// Замените эту строку на свою из MongoDB Atlas
	uri := "mongodb+srv://beathovenmozart:s1KDWX5%24@cluster0.zowoyqf.mongodb.net/product_checker?retryWrites=true&w=majority"

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB Atlas:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB Atlas:", err)
	}

	MongoClient = client
	MongoDB = client.Database("product_checker")
	log.Println("Successfully connected to MongoDB Atlas!")
}
func AddHistoryToMongo(username, productID, result string) error {
	collection := MongoDB.Collection("check_history")

	_, err := collection.InsertOne(context.Background(), models.ProductCheckHistory{
		Username:  username,
		ProductID: productID,
		Result:    result,
		CheckedAt: time.Now(),
	})

	return err
}

func GetHistoryByUsername(username string) ([]models.ProductCheckHistory, error) {
	collection := MongoDB.Collection("check_history")
	var history []models.ProductCheckHistory

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &history); err != nil {
		return nil, err
	}

	return history, nil
}

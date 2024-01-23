package app

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB configuration
var (
	mongoDBURL string
	MongoDB    string
	client     *mongo.Client
)

func InitMongoDB() {
	mongoDBURL = os.Getenv("MONGODB_URL")
	MongoDB = os.Getenv("MONGODB_DB")

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}

func GetMongoClient() *mongo.Client {
	return client
}

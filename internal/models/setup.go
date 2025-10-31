package models

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var mongoClient *mongo.Client
var collections map[string]*mongo.Collection

func ConnectDB() {
	uri := os.Getenv("MONGODB_URI")
	collections = make(map[string]*mongo.Collection)

	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	mongoClient = client
	collections["config"] = client.Database("signal").Collection("config")
	collections["error_log"] = client.Database("signal").Collection("error-logs")

	CreateConfig()
}

func DisconnectDB() {
	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		log.Fatal("Failed to disconnect from MongoDB: ", err)
	}
}

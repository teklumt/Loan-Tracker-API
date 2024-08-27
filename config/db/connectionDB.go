package db

import (
	"context"
	"log"
	"time"

	// "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var LoanCollection *mongo.Collection
var UserCollection *mongo.Collection
var LogCollection *mongo.Collection
func ConnectDB(connectionString string) {

    clientOptions := options.Client().ApplyURI(connectionString)

    client, err := mongo.NewClient(clientOptions)

    if err != nil {
        log.Fatalf("Error creating MongoDB client: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatalf("Error connecting to MongoDB: %v", err)
    }

    Client = client
    UserCollection = client.Database("loan_tracker_api").Collection("users")
    LoanCollection = client.Database("loan_tracker_api").Collection("loans")
    LogCollection = client.Database("loan_tracker_api").Collection("logs")
}

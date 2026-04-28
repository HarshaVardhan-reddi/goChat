package databases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// InitializeMongoDbClient reads the MONGODB connection URI from the .env file,
// connects to MongoDB, and verifies the connection with a ping.
//
// It returns a connected *mongo.Client, or an error if the .env file cannot be
// read, the MONGODB key is missing, or the connect/ping fails. On ping failure
// the client is disconnected before returning.
//
// Reference: https://www.mongodb.com/docs/drivers/go/current/


func InitializeMongoDbClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uriData, err := godotenv.Read()
	if err != nil {
		return nil, fmt.Errorf("read .env: %w", err)
	}

	uri := uriData["MONGODB"]
	if uri == "" {
		return nil, errors.New("MONGODB not set in .env")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	return client, nil
}

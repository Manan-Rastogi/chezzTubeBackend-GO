package db

import (
	"context"
	"log"
	"time"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func init() {
	connectDB()
}

func connectDB() {
	uri := configs.ENV.MongoDbUri
	if uri == "" {
		log.Fatalf("error reading mongo uri from env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Ensure context is canceled to avoid leaks

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(50))
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	// Check if the client connected successfully
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("error pinging MongoDB: %v", err)
	}

	// Connected successfully
	log.Println("Connected to MongoDB!")
}

func DisconnectDB() {
	if Client != nil {
		if err := Client.Disconnect(context.Background()); err != nil {
			log.Fatalf("error disconnecting from MongoDB: %v", err)
		}
	}
}

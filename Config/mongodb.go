package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(GetMongoURI()))
	if err != nil {
		log.Fatal("Gagal koneksi ke MongoDB:", err)
	}

	// Ping MongoDB
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Gagal ping MongoDB:", err)
	}

	log.Println("MongoDB connected successfully")

	// Return database
	return client.Database(GetMongoDatabase())
}

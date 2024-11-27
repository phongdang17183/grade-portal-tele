package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client


func ConnectMongoDB(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()


	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Không thể kết nối tới MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Không thể ping MongoDB: %v", err)
	}

	log.Println("Kết nối thành công tới MongoDB!")
	MongoClient = client
}



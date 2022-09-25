package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DatabaseConnection(uri string) *mongo.Client {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Panic(err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Println("Failed to ping the database")
		panic(err)
	}
	log.Println("Database successfully pinged ! ")
	return client

}

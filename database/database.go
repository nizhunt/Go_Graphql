package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionString string = "mongodb://localhost:27017"

type DB struct{
	client *mongo.Client
}

func Connect() *DB{
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err!= nil {
        log.Fatal(err)
    }
	err = client.Ping(ctx, readpref.Primary())
	if err!= nil {
        log.Fatal(err)
    }

	return &DB{
        client: client,
    }
}




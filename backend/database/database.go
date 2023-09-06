package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Connected successfully to MongoDB")
	return nil
}

func InitDb() (*mongo.Client, context.Context, context.CancelFunc, error) {
	client, ctx, cancel, err := Connect("mongodb://localhost:27017")
	if err != nil {
		return nil, nil, nil, err
	}
	defer func() {
		if err := Ping(client, ctx); err != nil {
			panic(err)
		}
	}()
	return client, ctx, cancel, nil
}

func SetupDatabase(client *mongo.Client) (*mongo.Collection, error) {
	goChatDB := client.Database("goChat")
	userCollection := goChatDB.Collection("userCollection")
	return userCollection, nil
}

package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection
var ctx = context.TODO()

type Crypto struct {
	Id      primitive.ObjectID
	Name    string
	Upvote  int64
	Dowvote int64
}

const uri = "mongodb://root:example@localhost:27017"

func Init() {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Collection = client.Database("tasker").Collection("tasks")

}

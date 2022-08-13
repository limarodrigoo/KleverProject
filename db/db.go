package db

import (
	"context"
	"log"

	pb "github.com/limarodrigoo/KleverProject/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection
var ctx = context.TODO()

type Crypto struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name" validate:"required"`
	Upvote   int64              `bson:"upvote" validate:"required"`
	Downvote int64              `bson:"downvote" validate:"required"`
}

const uri = "mongodb://root:example@localhost:27017"

var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("crypto").Collection("cryptos")

}

func CreateCryptoDb(crypto *pb.CryptoCreateReq) (*mongo.InsertOneResult, error) {
	insertedResult, err := collection.InsertOne(context.TODO(), crypto)

	if err != nil {
		return nil, err
	}

	return insertedResult, nil
}

func ListAllCryptos() (*mongo.Cursor, error) {
	opts := options.Find().SetSort(bson.D{{Key: "Upvote", Value: -1}})

	cursor, err := collection.Find(context.Background(), bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	return cursor, nil

}

func GetCryptoById(id primitive.ObjectID) *mongo.SingleResult {
	res := collection.FindOne(ctx, bson.M{"_id": id})

	return res

}

func UpvoteCryptById(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"Upvote": 1}}

	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}

	return nil
}

func DownvoteCryptById(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"Downvote": 1}}

	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}

	return nil
}

func DeleteCryptoById(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil
}

func GetCryptoByName(name string) bool {
	filter := bson.M{"name": name}

	res := collection.FindOne(ctx, filter)

	if res != nil {
		return true
	}
	return false
}

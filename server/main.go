package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/limarodrigoo/KleverProject/db"
	pb "github.com/limarodrigoo/KleverProject/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "Server port")
)

type server struct {
	pb.UnimplementedVotingServiceServer
}

func createCrypto(ctx context.Context) *mongo.InsertOneResult {
	crypto := bson.D{{Key: "Name", Value: "BTC"}, {Key: "Upvote", Value: 3}, {Key: "Downvote", Value: 1}}
	result, err := db.Collection.InsertOne(ctx, crypto)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.InsertedID)
	return result
}

func (s *server) CreateCrypto(ctx context.Context, in *pb.CryptoCreateReq) (*pb.CreateCryptoRes, error) {
	log.Printf("Received: %v", in.GetName())
	id := createCrypto(context.TODO())

	log.Printf("Id: %v", id.InsertedID)
	return &pb.CreateCryptoRes{Id: fmt.Sprint(id.InsertedID.(primitive.ObjectID))}, nil
}

func main() {
	flag.Parse()

	db.Init()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterVotingServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

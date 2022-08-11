package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/limarodrigoo/KleverProject/db"
	pb "github.com/limarodrigoo/KleverProject/proto"
	"github.com/limarodrigoo/KleverProject/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "Server port")
)

type server struct {
	pb.UnimplementedVotingServiceServer
}

func (s *server) CreateCrypto(ctx context.Context, in *pb.CryptoCreateReq) (*pb.CreateCryptoRes, error) {

	err := service.CheckValidation(in.GetName(), in.GetUpvote(), in.GetDownvote())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid input: %v", err))
	}

	crypto := bson.D{{Key: "Name", Value: in.GetName()}, {Key: "Upvote", Value: in.GetUpvote()}, {Key: "Downvote", Value: in.GetDownvote()}}

	result, err := db.Collection.InsertOne(context.TODO(), crypto)
	if err != nil {
		panic(err)
	}

	oid := result.InsertedID.(primitive.ObjectID)

	return &pb.CreateCryptoRes{Id: oid.Hex()}, nil
}

func (s *server) ListCryptos(in *pb.ListCryptosReq, stream pb.VotingService_ListCryptosServer) error {
	opts := options.Find().SetSort(bson.D{{Key: "Upvote", Value: -1}})

	cursor, err := db.Collection.Find(context.Background(), bson.D{}, opts)

	if err != nil {
		return status.Errorf(codes.NotFound, fmt.Sprintf("Ops, something went wrong: %v", err))
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.Background()) {
		result := db.Crypto{}
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		stream.Send(&pb.ListCryptosRes{
			Crypto: &pb.Crypto{
				Id:       result.Id.Hex(),
				Name:     result.Name,
				Upvote:   result.Upvote,
				Downvote: result.Downvote,
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknow cursor error: %v", err))
	}
	defer cursor.Close(context.TODO())
	return nil
}

func (s *server) GetCrypto(ctx context.Context, in *pb.GetCryptoReq) (*pb.Crypto, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := db.Collection.FindOne(ctx, bson.M{"_id": id})
	data := db.Crypto{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find crypto with ObjectId %s: %v", in.GetId(), err))
	}

	res := &pb.Crypto{
		Id:       id.Hex(),
		Name:     data.Name,
		Upvote:   data.Upvote,
		Downvote: data.Downvote,
	}

	return res, nil
}

func (s *server) UpvoteCrypto(ctx context.Context, in *pb.UpvoteCryptoReq) (*pb.UpvoteCryptoRes, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied crypto id to a ObjectId: %v", err),
		)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"Upvote": 1}}

	_, err = db.Collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find crypto with id %s: %v", in.GetId(), err))
	}

	return &pb.UpvoteCryptoRes{
		Success: true,
	}, nil
}

func (s *server) DownvoteCrypto(ctx context.Context, in *pb.DownvoteCryptoReq) (*pb.DownvoteCryptoRes, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied crypto id to a ObjectId: %v", err),
		)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"Downvote": 1}}

	_, err = db.Collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find crypto with id %s: %v", in.GetId(), err))
	}

	return &pb.DownvoteCryptoRes{
		Success: true,
	}, nil
}

func (s *server) DeleteCrypto(ctx context.Context, in *pb.DeleteCryptoReq) (*pb.DeleteCryptoRes, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	_, err = db.Collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cold not delete crypto with id: %s: %v", in.GetId(), err))
	}

	return &pb.DeleteCryptoRes{
		Success: true,
	}, nil
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

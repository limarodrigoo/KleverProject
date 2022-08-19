package helper

import (
	"context"
	"fmt"
	"log"

	"github.com/limarodrigoo/KleverProject/db"
	pb "github.com/limarodrigoo/KleverProject/proto"
	"github.com/limarodrigoo/KleverProject/server/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedVotingServiceServer
}

func (s *Server) CreateCrypto(ctx context.Context, in *pb.CryptoCreateReq) (*pb.CreateCryptoRes, error) {

	err := service.CheckValidation(in.GetName(), in.GetUpvote(), in.GetDownvote())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid input: %v", err))
	}

	crypto := pb.CryptoCreateReq{Name: in.GetName(), Upvote: in.GetUpvote(), Downvote: in.GetDownvote()}

	result, err := db.CreateCryptoDb(&crypto)

	oid := result.InsertedID.(primitive.ObjectID)

	return &pb.CreateCryptoRes{Id: oid.Hex()}, nil
}

func (s *Server) ListCryptos(in *pb.ListCryptosReq, stream pb.VotingService_ListCryptosServer) error {

	cursor, err := db.ListAllCryptos()

	if err != nil {
		return status.Errorf(codes.NotFound, fmt.Sprintf("Ops, something went wrong: %v", err))
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.Background()) {
		result := db.Crypto{}
		err := cursor.Decode(&result)
		if err != nil {
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
	return nil
}

func (s *Server) GetCrypto(ctx context.Context, in *pb.GetCryptoReq) (*pb.Crypto, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := db.GetCryptoById(id)

	data := db.Crypto{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find crypto with ObjectId %s", in.GetId()))
	}

	res := &pb.Crypto{
		Id:       id.Hex(),
		Name:     data.Name,
		Upvote:   data.Upvote,
		Downvote: data.Downvote,
	}

	return res, nil
}

func (s *Server) UpvoteCrypto(ctx context.Context, in *pb.UpvoteCryptoReq) (*pb.UpvoteCryptoRes, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied crypto id to a ObjectId: %v", err),
		)
	}

	err = db.UpvoteCryptById(id)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find crypto with id %s: %v", in.GetId(), err))
	}

	return &pb.UpvoteCryptoRes{
		Success: true,
	}, nil
}

func (s *Server) DownvoteCrypto(ctx context.Context, in *pb.DownvoteCryptoReq) (*pb.DownvoteCryptoRes, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied crypto id to a ObjectId: %v", err),
		)
	}

	err = db.DownvoteCryptById(id)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find crypto with id %s: %v", in.GetId(), err))
	}

	return &pb.DownvoteCryptoRes{
		Success: true,
	}, nil
}

func (s *Server) DeleteCrypto(ctx context.Context, in *pb.DeleteCryptoReq) (*pb.DeleteCryptoRes, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	err = db.DeleteCryptoById(id)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cold not delete crypto with id: %s: %v", in.GetId(), err))
	}

	return &pb.DeleteCryptoRes{
		Success: true,
	}, nil
}

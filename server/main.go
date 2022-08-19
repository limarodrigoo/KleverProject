package main

import (
	"log"
	"net"

	pb "github.com/limarodrigoo/KleverProject/proto"
	"github.com/limarodrigoo/KleverProject/server/helper"
	"google.golang.org/grpc"
)

var (
	port = ":50051"
)

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterVotingServiceServer(s, &helper.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

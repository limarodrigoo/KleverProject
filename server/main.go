package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/limarodrigoo/KleverProject/db"
	pb "github.com/limarodrigoo/KleverProject/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "Server port")
)

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

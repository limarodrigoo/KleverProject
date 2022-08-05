package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/limarodrigoo/KleverProject/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", "BTC", "Crypto name")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewVotingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateCrypto(ctx, &pb.CryptoCreateReq{Name: *name})
	if err != nil {
		log.Fatalf("could not create: %v", err)
	}
	log.Printf("Crypto id: %v\n", r.GetId())
	log.Printf("Crypto name: %v\n", *name)
}

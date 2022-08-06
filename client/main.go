package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/limarodrigoo/KleverProject/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr   = flag.String("addr", "localhost:50051", "the address to connect to")
	client pb.VotingServiceClient
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client = pb.NewVotingServiceClient(conn)

	router := gin.Default()

	router.POST("/crypto", createCrypto)
	router.GET("/crypto/:id", getCrypto)
	router.PUT("/upvote/:id", upvoteCrypto)
	router.PUT("/downvote/:id", downvoteCrypto)
	router.GET("/cryptos", listCryptos)
	router.DELETE("/crypto/:id", deleteCrypto)
	log.Fatal(router.Run(":8080"))

}

func createCrypto(ctx *gin.Context) {
	crypto := pb.CryptoCreateReq{}
	err := ctx.ShouldBindJSON(&crypto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(&crypto)

	res, err := client.CreateCrypto(ctx, &crypto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(res),
		})
	}
}

func upvoteCrypto(ctx *gin.Context) {
	id := ctx.Param("id")

	obj := pb.UpvoteCryptoReq{Id: id}

	res, err := client.UpvoteCrypto(ctx, &obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func downvoteCrypto(ctx *gin.Context) {
	id := ctx.Param("id")

	obj := pb.DownvoteCryptoReq{Id: id}

	res, err := client.DownvoteCrypto(ctx, &obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func getCrypto(ctx *gin.Context) {
	id := ctx.Param("id")

	obj := pb.GetCryptoReq{Id: id}

	res, err := client.GetCrypto(ctx, &obj)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func deleteCrypto(ctx *gin.Context) {
	id := ctx.Param("id")

	obj := pb.DeleteCryptoReq{Id: id}

	res, err := client.DeleteCrypto(ctx, &obj)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func listCryptos(ctx *gin.Context) {
	obj := pb.ListCryptosReq{}

	stream, err := client.ListCryptos(ctx, &obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(res.Crypto.GetId())
		ctx.JSON(http.StatusOK, res.GetCrypto())
	}
}

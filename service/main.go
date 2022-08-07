package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckValidation(name string, upvote int64, downvote int64) error {
	if name == "" {
		return status.Errorf(codes.InvalidArgument, "Name is required!")
	}

	if downvote != 0 || upvote != 0 {
		return status.Errorf(codes.InvalidArgument, "Cryptos must be initialized with 0 votes")
	}

	return nil
}

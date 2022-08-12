package service_test

import (
	"testing"

	"github.com/limarodrigoo/KleverProject/server/service"
)

func TestValidation(t *testing.T) {
	test := []struct {
		name    string
		upvote  int64
		dowvote int64
		wanted  string
	}{
		{
			name:    "BTC",
			upvote:  0,
			dowvote: -1,
			wanted:  "rpc error: code = InvalidArgument desc = Cryptos must be initialized with 0 votes",
		},
		{
			name:    "BTC",
			upvote:  2,
			dowvote: 0,
			wanted:  "rpc error: code = InvalidArgument desc = Cryptos must be initialized with 0 votes",
		},
		{
			name:    "",
			upvote:  0,
			dowvote: 0,
			wanted:  "rpc error: code = InvalidArgument desc = Name is required!",
		},
	}

	for _, tt := range test {
		err := service.CheckValidation(tt.name, tt.upvote, tt.dowvote)
		if err.Error() != tt.wanted {
			t.Errorf("Expected: %v, Have: %v", tt.wanted, err.Error())
		}
	}
}

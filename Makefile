build:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/voting.proto
	@go mod download
run:
	@go run server/main.go
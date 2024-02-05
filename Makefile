build:
	go build ./cmd/server

run:
	go run ./cmd/server

swagger:
	swag init -g ./cmd/server/main.go

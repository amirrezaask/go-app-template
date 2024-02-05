lint:
	golangci-lint run

build:  
	go build ./cmd/server
	lint

run:
	go run ./cmd/server

swagger:
	swag init -g ./cmd/server/main.go


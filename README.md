# Go App Template


## How to use
1. First replace all {{template}} with your module path.
2. Enjoy.

## Libraries

- [Echo Web Server](https://echo.labstack.com/)
- [Swagger](https://github.com/swaggo/swag)
- [MySQL Driver](https://github.com/go-sql-driver/mysql)
- [Redis Driver](https://github.com/redis/go-redis)
- [RabbitMQ Driver](https://github.com/rabbitmq/amqp091-go)
- [SQLX Database Helper](https://github.com/jmoiron/sqlx)
- [OpenTelemetryGo](https://github.com/open-telemetry/opentelemetry-go)
- [go-elasticsearch](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/overview.html)
- [Minio](https://github.com/minio/minio-go)

## System Requirements

- [Golang toolchain](https://go.dev/dl/)


- swag, goimports, golangci-lint 
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
go install github.com/swaggo/swag/cmd/swag@latest
go install golang.org/x/tools/cmd/goimports
go install github.com/cosmtrek/air@latest
```

## Running
```
air
```

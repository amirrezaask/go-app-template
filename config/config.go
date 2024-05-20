package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func getEnv(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}

	return value
}

var (
	DEBUG                      bool
	TRACING_SERVICE_NAME       string
	LISTEN_ADDR                string
	LOG_LEVEL                  slog.Level
	DATABASE_CONNECTION_STRING string
	REDIS_HOST                 string
	REDIS_DB                   int
	ELASTIC_SEARCH_ADDRESSES   []string
	ELASTIC_SEARCH_USERNAME    string
	ELASTIC_SEARCH_PASSWORD    string
	MINIO_BUCKET_NAME          string
	MINIO_ENDPOINT             string
	MINIO_ACCESS_ID            string
	MINIO_SECRET_ACCESS_ID     string
	JWT_SECRET_KEY             []byte
)

func initVariables() {
	var err error

	ll := getEnv("LOG_LEVEL", "debug")
	if ll != "" {
		switch ll {
		case "debug":
			LOG_LEVEL = slog.LevelDebug
		case "info":
			LOG_LEVEL = slog.LevelInfo
		case "warn":
			LOG_LEVEL = slog.LevelWarn
		case "error":
			LOG_LEVEL = slog.LevelError
		default:
			LOG_LEVEL = slog.LevelInfo
		}
	}
	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		panic(err)
	}
	REDIS_DB = db
	TRACING_SERVICE_NAME = getEnv("TRACING_SERVICE_NAME", "template")
	LISTEN_ADDR = getEnv("LISTEN_ADDR", ":8080")
	DATABASE_CONNECTION_STRING = getEnv("DATABASE_CONNECTION_STRING", "")
	REDIS_HOST = getEnv("REDIS_HOST", "localhost:6379")
	ELASTIC_SEARCH_ADDRESSES = strings.Split(getEnv("ELASTIC_SEARCH_ADDRESSES", ""), ",")
	ELASTIC_SEARCH_USERNAME = getEnv("ELASTIC_SEARCH_USERNAME", "")
	ELASTIC_SEARCH_PASSWORD = getEnv("ELASTIC_SEARCH_PASSWORD", "")
	MINIO_BUCKET_NAME = getEnv("MINIO_BUCKET_NAME", "")
	MINIO_ENDPOINT = getEnv("MINIO_ENDPOINT", "")
	MINIO_ACCESS_ID = getEnv("MINIO_ACCESS_ID", "")
	MINIO_SECRET_ACCESS_ID = getEnv("MINIO_SECRET_ACCESS_ID", "")
	JWT_SECRET_KEY = []byte(getEnv("JWT_SECRET_KEY", ""))
}

func init() {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	initVariables()
}

package main

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

func getEnv(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}

	return value
}

var (
	TRACING_SERVICE_NAME       string
	LISTEN_ADDR                string
	LOG_LEVEL                  zerolog.Level
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
	LOG_LEVEL, err = zerolog.ParseLevel(getEnv("LOG_LEVEL", "debug"))
	if err != nil {
		panic(err)
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

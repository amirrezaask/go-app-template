package main

import (
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
	TRACING_SERVICE_NAME       = getEnv("TRACING_SERVICE_NAME", "template")
	LISTEN_ADDR                = getEnv("LISTEN_ADDR", ":8080")
	LOG_LEVEL                  zerolog.Level
	DATABASE_CONNECTION_STRING = getEnv("DATABASE_CONNECTION_STRING", "mysql:mysql@tcp(localhost:3306)/database")
	REDIS_HOST                 = getEnv("REDIS_HOST", "localhost:6379")
	REDIS_DB                   int
	ELASTIC_SEARCH_ADDRESSES   = strings.Split(getEnv("ELASTIC_SEARCH_ADDRESSES", ""), ",")
	ELASTIC_SEARCH_USERNAME    = getEnv("ELASTIC_SEARCH_USERNAME", "")
	ELASTIC_SEARCH_PASSWORD    = getEnv("ELASTIC_SEARCH_PASSWORD", "")
)

func init() {
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

}

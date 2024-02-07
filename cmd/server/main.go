package main

import (
	"context"

	"gitlab.snappcloud.io/doctor/backend/template/api"
	"gitlab.snappcloud.io/doctor/backend/template/logger"
	"gitlab.snappcloud.io/doctor/backend/template/storage"
	"gitlab.snappcloud.io/doctor/backend/template/tracing"
)

// @title Template API
// @version 1.0

// @contact.name Doctor Backend

// @host snapp.doctor
// @BasePath /
func main() {
	_ = context.TODO()
	l := logger.New(LOG_LEVEL)

	shutdown, err := tracing.Init()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := shutdown(context.Background()); err != nil {
			l.Error().Err(err).Send()
		}
	}()

	db, err := storage.NewMySQL(DATABASE_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}

	redis, err := storage.NewRedis(context.Background(), REDIS_HOST, REDIS_DB)
	if err != nil {
		panic(err)
	}

	e := api.NewAPIServer(db, redis, l)

	if err := e.Start(LISTEN_ADDR); err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"log/slog"
	"os"

	"{{template}}/api"
	"{{template}}/config"
	"{{template}}/storage"
	"{{template}}/tracing"
)

// @title Go Service Template API
// @version 1.0

// @host app.awesome
// @BasePath /
func main() {
	globalLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     config.LOG_LEVEL,
	}))

	if !config.DEBUG {
		globalLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     config.LOG_LEVEL,
		}))
	}
	slog.SetDefault(globalLogger)
	slog.SetLogLoggerLevel(config.LOG_LEVEL)

	shutdown, err := tracing.Init()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := shutdown(context.Background()); err != nil {
			slog.Error("error in shutting down tracing", "err", err)
		}
	}()

	db, err := storage.NewMySQL(config.DATABASE_CONNECTION_STRING)
	if err != nil {
		panic(err)
	}

	redis, err := storage.NewRedis(context.Background(), config.REDIS_HOST, config.REDIS_DB)
	if err != nil {
		panic(err)
	}

	e := api.NewAPIServer(config.TRACING_SERVICE_NAME, db, redis)

	if err := e.Start(config.LISTEN_ADDR); err != nil {
		panic(err)
	}
}

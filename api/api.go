package api

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	_ "gitlab.snappcloud.io/doctor/backend/template/docs"
	"gitlab.snappcloud.io/doctor/backend/template/logger"
	"gitlab.snappcloud.io/doctor/backend/template/storage"
)

func NewAPIServer(db *storage.MySQL, redis *storage.Redis, l logger.Logger) *echo.Echo {
	e := echo.New()

	e.Use(echoprometheus.NewMiddleware("app_template"))

	e.GET("/metrics", echoprometheus.NewHandler())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	apiG := e.Group("/api")
	v1G := e.Group("/v1")

	_ = apiG
	_ = v1G

	return e
}

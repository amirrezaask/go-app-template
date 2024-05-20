package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	echoPrometheus "github.com/globocom/echo-prometheus"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "gitlab.snappcloud.io/doctor/backend/template/docs"
	"gitlab.snappcloud.io/doctor/backend/template/storage"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// jsonIterSerializer implements JSON encoding using encoding/json.
type jsonIterSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d jsonIterSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := jsoniter.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d jsonIterSerializer) Deserialize(c echo.Context, i interface{}) error {
	err := jsoniter.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}

func setTracingSpanMiddleware(service string) echo.MiddlewareFunc {
	return otelecho.Middleware(service)
}

func NewAPIServer(tracingServiceName string, db *storage.MySQL, redis *storage.Redis) *echo.Echo {
	e := echo.New()
	e.JSONSerializer = &jsonIterSerializer{}

	collector := echoPrometheus.MetricsMiddlewareWithConfig(echoPrometheus.Config{
		Namespace: "app_template",
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	apiG := e.Group("/api", setTracingSpanMiddleware(tracingServiceName), collector)
	v1G := e.Group("/v1")

	_ = apiG
	_ = v1G

	return e
}

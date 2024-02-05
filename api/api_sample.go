package api

import (
	"github.com/labstack/echo/v4"

	"gitlab.snappcloud.io/doctor/backend/template/logger"
)

type SampleHandler struct {
	l logger.Logger
}

type sampleIndexRequest struct {
}

type sampleIndexResponse struct {
}

// Index ...
// @Summary      Sample Index
// @Accept       json
// @Produce      json
// @Param        req   body      sampleIndexRequest  true  "Index request"
// @Success      200  {object}   sampleIndexResponse
// @Router       /api/v1 [get]
func (s *SampleHandler) Index(c echo.Context) error {
	var req sampleIndexRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, sampleIndexResponse{})
}

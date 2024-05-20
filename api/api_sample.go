package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SampleHandler struct {
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
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, sampleIndexResponse{})
}

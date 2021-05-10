package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) getTokens(c echo.Context) error {
	resp := make([]*TokenResponse, 0)
	return c.JSON(http.StatusOK, resp)
}


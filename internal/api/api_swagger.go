package api

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"os"
)

func (s *Server) getSwagger(c echo.Context) error {
	f, err := os.Open("./api/swagger/marketplace.swagger.json")
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "application/json", data)
}

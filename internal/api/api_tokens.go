package api

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"net/http"
	"strconv"
)

func (s *Server) getTokens(c echo.Context) error {
	ctx := context.Background()

	symbol := c.FormValue("symbol")
	address := c.FormValue("address")
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)
	fltr := &datastore.TokensFilter{}
	if symbol != "" {
		fltr.Symbol = pointer.ToString(symbol)
	}
	if address != "" {
		fltr.Address = pointer.ToString(address)
	}

	tokens, err := s.ds.Tokens.List(ctx, fltr, limitOpts)
	if err != nil {
		return err
	}

	resp := toTokensResponse(tokens)
	return c.JSON(http.StatusOK, resp)
}


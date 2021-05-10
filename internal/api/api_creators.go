package api

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) GetCreators(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)

	fltr := &datastore.AccountsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "created_at",
			IsAsc: true,
		},
	}

	q := strings.TrimSpace(c.FormValue("q"))
	if len(q) > 0 {
		fltr.Query = pointer.ToString(q)
	}

	ctx := context.Background()
	creators, err := s.ds.Accounts.List(ctx, fltr, limitOpts)
	if err != nil {
		return err
	}

	tc, _ := s.ds.Accounts.Count(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     *limitOpts.Offset,
		Limit:      *limitOpts.Limit,
	}

	resp := toAccountsResponse(creators, countResp)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) GetCreator(c echo.Context) error {
	creatorID, _ := strconv.ParseInt(c.Param("creator_id"), 10, 64)
	if creatorID == 0 {
		return echo.ErrNotFound
	}

	ctx := context.Background()
	creator, err := s.ds.Accounts.GetByID(ctx, creatorID)
	if err != nil {
		if err == datastore.ErrAccountNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	resp := toAccountResponse(creator)
	return c.JSON(http.StatusOK, resp)
}

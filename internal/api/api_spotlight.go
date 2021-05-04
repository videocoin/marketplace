package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"net/http"
	"strconv"
)

func (s *Server) getSpotlightFeaturedArts(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)

	fltr := &datastore.ArtsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "id",
			IsAsc: true,
		},
	}

	ctx := context.Background()
	arts, err := s.ds.GetArtsList(ctx, fltr, limitOpts)
	if err != nil {
		return err
	}

	tc, _ := s.ds.GetArtsListCount(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     *limitOpts.Offset,
		Limit:      *limitOpts.Limit,
	}

	resp := toArtsResponse(arts, countResp)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getSpotlightLiveArts(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)

	fltr := &datastore.ArtsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "name",
			IsAsc: false,
		},
	}

	ctx := context.Background()
	arts, err := s.ds.GetArtsList(ctx, fltr, limitOpts)
	if err != nil {
		return err
	}

	tc, _ := s.ds.GetArtsListCount(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     *limitOpts.Offset,
		Limit:      *limitOpts.Limit,
	}

	resp := toArtsResponse(arts, countResp)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getSpotlightFeaturedCreators(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)

	fltr := &datastore.AccountsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "created_at",
			IsAsc: true,
		},
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

	resp := toCreatorsResponse(creators, countResp)
	return c.JSON(http.StatusOK, resp)
}

package api

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) getActivity(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)

	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)
	groupID := strings.TrimSpace(c.FormValue("group_id"))

	ctx := context.Background()
	fltr := &datastore.ActivityFilter{
		CreatedByID: pointer.ToInt64(account.ID),
		Sort: &datastore.SortOption{
			Field: "created_at",
			IsAsc: false,
		},
	}
	if groupID != "" {
		fltr.GroupID = pointer.ToString(groupID)
	}

	items, err := s.ds.Activity.List(ctx, fltr, limitOpts)
	if err != nil {
		return err
	}

	err = s.ds.JoinAssetToActivity(ctx, items)
	if err != nil {
		return err
	}

	err = s.ds.JoinOrderToActivity(ctx, items)
	if err != nil {
		return err
	}

	tc, _ := s.ds.Activity.Count(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     *limitOpts.Offset,
		Limit:      *limitOpts.Limit,
	}

	resp := toActivityResponse(items, countResp)
	return c.JSON(http.StatusOK, resp)
}

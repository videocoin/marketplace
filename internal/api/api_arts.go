package api

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	pkgyt "github.com/videocoin/marketplace/pkg/youtube"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) createArt(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)

	art := &model.Art{
		CreatedByID: account.ID,
		Account:     account,
	}

	req := new(CreateArtRequest)
	err := c.Bind(req)
	if err != nil {
		return echo.ErrBadRequest
	}

	artName := strings.TrimSpace(req.Name)
	if artName == "" {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "missing name"})
	}
	art.Name = artName

	if req.AssetID == 0 {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "asset not found"})
	}

	ctx := context.Background()
	asset, err := s.ds.Assets.GetByID(ctx, req.AssetID)
	if err != nil {
		if err == datastore.ErrAssetNotFound {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "asset not found"})
		}

		return err
	}

	if asset.CreatedByID != account.ID {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "asset not found"})
	}

	art.Asset = asset
	art.AssetID = asset.ID

	if req.Description != nil {
		artDesc := strings.TrimSpace(*req.Description)
		if artDesc != "" {
			art.Desc = dbr.NewNullString(artDesc)
		}
	}

	if req.YoutubeLink != nil {
		ytLink, err := pkgyt.ValidateVideoURL(*req.YoutubeLink)
		if err != nil {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": "wrong youtube link"})
		}

		art.YTLink = dbr.NewNullString(ytLink)
	}

	err = s.ds.Arts.Create(ctx, art)
	if err != nil {
		return err
	}

	resp := toArtResponse(art)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getArts(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)
	fltr := &datastore.ArtsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "created_at",
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

func (s *Server) getArtsByCreator(c echo.Context) error {
	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)

	creatorID, _ := strconv.ParseInt(c.Param("creator_id"), 10, 64)
	if creatorID == 0 {
		return echo.ErrNotFound
	}

	fltr := &datastore.ArtsFilter{
		CreatedByID: pointer.ToInt64(creatorID),
		Sort: &datastore.DatastoreSort{
			Field: "created_at",
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

func (s *Server) getArt(c echo.Context) error {
	ctx := context.Background()

	artID, _ := strconv.ParseInt(c.Param("art_id"), 10, 64)
	if artID == 0 {
		return echo.ErrNotFound
	}

	art, err := s.ds.Arts.GetByID(ctx, artID)
	if err != nil {
		if err == datastore.ErrArtNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	asset, err := s.ds.Assets.GetByID(ctx, art.AssetID)
	if err != nil {
		return err
	}

	art.Asset = asset

	resp := toArtResponse(art)
	return c.JSON(http.StatusOK, resp)
}

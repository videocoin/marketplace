package marketplace

import (
	"context"
	"errors"
	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/api/rpc"
	v1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/rpcauth"
	pkgyt "github.com/videocoin/marketplace/pkg/youtube"
	"strings"
)

func (s *Server) CreateArt(ctx context.Context, req *v1.CreateArtRequest) (*v1.ArtResponse, error) {
	account, _ := rpcauth.AccountFromContext(ctx)
	art := &model.Art{
		CreatedByID: account.ID,
		Account:     account,
	}

	artName := strings.TrimSpace(req.Name)
	if artName == "" {
		return nil, rpc.NewRpcValidationError(&rpc.ValidationError{
			Field:   "name",
			Message: "missing name",
		})
	}
	art.Name = artName

	if req.AssetId == 0 {
		return nil, rpc.NewRpcSimpleValidationError(errors.New("asset not found"))
	}

	asset, err := s.ds.Assets.GetByID(ctx, req.AssetId)
	if err != nil {
		if err == datastore.ErrAssetNotFound {
			return nil, rpc.NewRpcSimpleValidationError(errors.New("asset not found"))
		}

		return nil, rpc.NewRpcInternalError(err)
	}

	if asset.CreatedByID != account.ID {
		return nil, rpc.NewRpcSimpleValidationError(errors.New("asset not found"))
	}

	art.Asset = asset
	art.AssetID = asset.ID

	if req.Description != nil {
		artDesc := strings.TrimSpace(req.Description.Value)
		if artDesc != "" {
			art.Desc = dbr.NewNullString(artDesc)
		}
	}

	if req.YoutubeLink != nil {
		ytLink, err := pkgyt.ValidateVideoURL(req.YoutubeLink.Value)
		if err != nil {
			return nil, rpc.NewRpcValidationError(&rpc.ValidationError{
				Field:   "youtube_link",
				Message: "wrong youtube link",
			})
		}

		art.YTLink = dbr.NewNullString(ytLink)
	}

	err = s.ds.Arts.Create(ctx, art)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	resp := toArtResponse(art)
	return resp, nil
}

func (s *Server) GetArts(ctx context.Context, req *v1.ArtsRequest) (*v1.ArtsResponse, error) {
	limitOpts := datastore.NewLimitOpts(req.Offset, req.Limit)
	fltr := &datastore.ArtsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "created_at",
			IsAsc: false,
		},
	}
	arts, err := s.ds.GetArtsList(ctx, fltr, limitOpts)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	tc, _ := s.ds.GetArtsListCount(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     req.Offset,
		Limit:      req.Limit,
	}

	resp := toArtsResponse(arts, countResp)
	return resp, nil
}

func (s *Server) GetFeaturedArts(ctx context.Context, req *v1.ArtsRequest) (*v1.ArtsResponse, error) {
	limitOpts := datastore.NewLimitOpts(req.Offset, req.Limit)
	fltr := &datastore.ArtsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "id",
			IsAsc: true,
		},
	}
	arts, err := s.ds.GetArtsList(ctx, fltr, limitOpts)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	tc, _ := s.ds.GetArtsListCount(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     req.Offset,
		Limit:      req.Limit,
	}

	resp := toArtsResponse(arts, countResp)
	return resp, nil
}

func (s *Server) GetLiveArts(ctx context.Context, req *v1.ArtsRequest) (*v1.ArtsResponse, error) {
	limitOpts := datastore.NewLimitOpts(req.Offset, req.Limit)
	fltr := &datastore.ArtsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "name",
			IsAsc: true,
		},
	}
	arts, err := s.ds.GetArtsList(ctx, fltr, limitOpts)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	tc, _ := s.ds.GetArtsListCount(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     req.Offset,
		Limit:      req.Limit,
	}

	resp := toArtsResponse(arts, countResp)
	return resp, nil
}

func (s *Server) GetArtsByCreator(ctx context.Context, req *v1.ArtsByCreatorRequest) (*v1.ArtsResponse, error) {
	limitOpts := datastore.NewLimitOpts(req.Offset, req.Limit)
	fltr := &datastore.ArtsFilter{
		CreatedByID: pointer.ToInt64(req.CreatorId),
		Sort: &datastore.DatastoreSort{
			Field: "created_at",
			IsAsc: false,
		},
	}
	arts, err := s.ds.GetArtsList(ctx, fltr, limitOpts)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	tc, _ := s.ds.GetArtsListCount(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     req.Offset,
		Limit:      req.Limit,
	}

	resp := toArtsResponse(arts, countResp)
	return resp, nil
}

func (s *Server) GetArt(ctx context.Context, req *v1.ArtRequest) (*v1.ArtResponse, error) {
	art, err := s.ds.Arts.GetByID(ctx, req.Id)
	if err != nil {
		if err == datastore.ErrArtNotFound {
			return nil, rpc.ErrRpcNotFound
		}
		return nil, rpc.NewRpcInternalError(err)
	}

	asset, err := s.ds.Assets.GetByID(ctx, art.AssetID)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	art.Asset = asset

	resp := toArtResponse(art)
	return resp, nil
}

package marketplace

import (
	"context"
	"github.com/videocoin/marketplace/api/rpc"
	v1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/datastore"
)

func (s *Server) GetSpotlightFeaturedArts(ctx context.Context, req *v1.ArtsRequest) (*v1.ArtsResponse, error) {
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

func (s *Server) GetSpotlightLiveArts(ctx context.Context, req *v1.ArtsRequest) (*v1.ArtsResponse, error) {
	limitOpts := datastore.NewLimitOpts(req.Offset, req.Limit)
	fltr := &datastore.ArtsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "name",
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

func (s *Server) GetSpotlightFeaturedCreators(ctx context.Context, req *v1.CreatorsRequest) (*v1.CreatorsResponse, error) {
	limitOpts := datastore.NewLimitOpts(req.Offset, req.Limit)
	fltr := &datastore.AccountsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "created_at",
			IsAsc: true,
		},
	}
	creators, err := s.ds.Accounts.List(ctx, fltr, limitOpts)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	tc, _ := s.ds.Accounts.Count(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     req.Offset,
		Limit:      req.Limit,
	}

	resp := toCreatorsResponse(creators, countResp)
	return resp, nil
}

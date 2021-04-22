package marketplace

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/videocoin/marketplace/api/rpc"
	v1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/rpcauth"
)

func (s *Server) GetMyArts(ctx context.Context, req *v1.ArtsRequest) (*v1.ArtsResponse, error) {
	account, _ := rpcauth.AccountFromContext(ctx)

	limitOpts := datastore.NewLimitOpts(req.Offset, req.Limit)
	fltr := &datastore.ArtsFilter{
		CreatedByID: pointer.ToInt64(account.ID),
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

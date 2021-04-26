package marketplace

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/videocoin/marketplace/api/rpc"
	v1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/datastore"
	"strings"
)

func (s *Server) GetCreators(ctx context.Context, req *v1.CreatorsRequest) (*v1.CreatorsResponse, error) {
	limitOpts := datastore.NewLimitOpts(req.Offset, req.Limit)
	fltr := &datastore.AccountsFilter{
		Sort: &datastore.DatastoreSort{
			Field: "created_at",
			IsAsc: true,
		},
	}
	q := strings.TrimSpace(req.Q)
	if len(q) > 0 {
		fltr.Query = pointer.ToString(q)
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

func (s *Server) GetCreator(ctx context.Context, req *v1.CreatorRequest) (*v1.CreatorResponse, error) {
	creator, err := s.ds.Accounts.GetByID(ctx, req.Id)
	if err != nil {
		if err == datastore.ErrAccountNotFound {
			return nil, rpc.ErrRpcNotFound
		}
		return nil, rpc.NewRpcInternalError(err)
	}

	resp := toCreatorResponse(creator)
	return resp, nil
}

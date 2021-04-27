package marketplace

import (
	"github.com/gogo/protobuf/types"
	v1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/model"
)

func toAssetResponse(asset *model.Asset) *v1.AssetResponse {
	resp := &v1.AssetResponse{
		Id:          asset.ID,
		ContentType: asset.ContentType,
		Status:      asset.Status,
	}

	if asset.Status == v1.AssetStatusReady {
		resp.Url = &types.StringValue{Value: asset.GetURL()}
		resp.ThumbnailUrl = &types.StringValue{Value: asset.GetThumbnailURL()}
		resp.PlaybackUrl = &types.StringValue{Value: asset.GetPlaybackURL()}
	}

	return resp
}

func toArtResponse(art *model.Art) *v1.ArtResponse {
	resp := &v1.ArtResponse{
		Id:   art.ID,
		Name: art.Name,
	}

	if art.Desc.Valid {
		resp.Description = art.Desc.String
	}

	if art.Asset != nil {
		resp.Asset = toAssetResponse(art.Asset)
	}

	if art.Account != nil {
		resp.Creator = toCreatorResponse(art.Account)
	}

	return resp
}

func toArtsResponse(arts []*model.Art, count *ItemsCountResponse) *v1.ArtsResponse {
	resp := &v1.ArtsResponse{
		Items: []*v1.ArtResponse{},
	}

	for _, art := range arts {
		resp.Items = append(resp.Items, toArtResponse(art))
	}

	resp.Count = int64(len(resp.Items))
	if count != nil {
		resp.TotalCount = count.TotalCount
		resp.Prev = resp.Count > 0 && count.Offset > 0
		resp.Next = resp.Count > 0 && resp.TotalCount > (resp.Count+int64(count.Offset))
	}

	return resp
}

func toCreatorResponse(creator *model.Account) *v1.CreatorResponse {
	resp := &v1.CreatorResponse{
		Id:      creator.ID,
		Address: creator.Address,
	}

	if creator.Username.Valid {
		resp.Username = &types.StringValue{Value: creator.Username.String}
	}

	if creator.Name.Valid {
		resp.Name = &types.StringValue{Value: creator.Name.String}
	}

	return resp
}

func toCreatorsResponse(creators []*model.Account, count *ItemsCountResponse) *v1.CreatorsResponse {
	resp := &v1.CreatorsResponse{
		Items: []*v1.CreatorResponse{},
	}

	for _, creator := range creators {
		resp.Items = append(resp.Items, toCreatorResponse(creator))
	}

	resp.Count = int64(len(resp.Items))
	if count != nil {
		resp.TotalCount = count.TotalCount
		resp.Prev = resp.Count > 0 && count.Offset > 0
		resp.Next = resp.Count > 0 && resp.TotalCount > (resp.Count+int64(count.Offset))
	}

	return resp
}

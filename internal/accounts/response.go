package accounts

import (
	"github.com/gogo/protobuf/types"
	v1 "github.com/videocoin/marketplace/api/v1/accounts"
	"github.com/videocoin/marketplace/internal/model"
)

func toNonceResponse(account *model.Account) *v1.NonceResponse {
	return &v1.NonceResponse{
		Nonce: account.Nonce.String,
	}
}

func toRegisterResponse(account *model.Account) *v1.RegisterResponse {
	return &v1.RegisterResponse{
		Address: account.Address,
		Nonce:   account.Nonce.String,
	}
}

func toAccountResponse(account *model.Account) *v1.AccountResponse {
	resp := &v1.AccountResponse{
		Address:  account.Address,
	}

	if account.Username.Valid {
		resp.Username = &types.StringValue{
			Value: account.Username.String,
		}
	}

	if account.Name.Valid {
		resp.Name = &types.StringValue{
			Value: account.Name.String,
		}
	}

	if account.ImageURL.Valid {
		resp.ImageUrl = &types.StringValue{
			Value: account.ImageURL.String,
		}
	}

	return resp
}

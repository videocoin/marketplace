package accounts

import (
	"context"
	"errors"
	"github.com/AlekSi/pointer"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/api/rpc"
	v1 "github.com/videocoin/marketplace/api/v1/accounts"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/rpcauth"
	"github.com/videocoin/marketplace/pkg/grpcutil"
	"strings"
)

type Server struct {
	logger     *logrus.Entry
	validator  *grpcutil.RequestValidator
	authSecret string
	ds         *datastore.Datastore
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	s := &Server{
		logger: ctxlogrus.Extract(ctx).WithField("system", "users"),
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return rpcauth.Auth(ctx, fullMethodName, s.authSecret, s.ds)
}

func (s *Server) GetNonce(ctx context.Context, req *v1.GetNonceRequest) (*v1.NonceResponse, error) {
	address := strings.ToLower(req.Address)
	if !ethcommon.IsHexAddress(address) {
		return nil, rpc.ErrRpcNotFound
	}

	account, err := s.ds.Accounts.GetByAddress(ctx, address)
	if err != nil {
		if err == datastore.ErrAccountNotFound {
			return nil, rpc.ErrRpcNotFound
		}
		return nil, rpc.NewRpcInternalError(err)
	}

	resp := toNonceResponse(account)
	return resp, nil
}

func (s *Server) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	address := strings.ToLower(req.Address)
	if !ethcommon.IsHexAddress(address) {
		return nil, rpc.NewRpcSimpleValidationError(ErrInvalidAddress)
	}

	_, err := s.ds.Accounts.GetByAddress(ctx, address)
	if err == nil {
		return nil, rpc.NewRpcSimpleValidationError(ErrAddressAlreadyRegistered)
	}

	if err != datastore.ErrAccountNotFound {
		return nil, rpc.NewRpcInternalError(err)
	}

	account := &model.Account{Address: address}
	err = s.ds.Accounts.Create(ctx, account)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	resp := toRegisterResponse(account)

	return resp, nil
}

func (s *Server) Auth(ctx context.Context, req *v1.AuthRequest) (*v1.AuthResponse, error) {
	address := strings.ToLower(req.Address)
	if !ethcommon.IsHexAddress(address) {
		return nil, rpc.NewRpcSimpleValidationError(ErrInvalidAddress)
	}

	logger := s.logger.
		WithField("address", address).
		WithField("signature", req.Signature)

	account, err := s.ds.Accounts.GetByAddress(ctx, address)
	if err != nil {
		if err == datastore.ErrAccountNotFound {
			return nil, rpc.NewRpcSimpleValidationError(ErrAddressNotRegistered)
		}
		return nil, rpc.NewRpcInternalError(err)
	}

	pkHex, err := verifySignature(req.Signature, account.Nonce.String, account.Address)
	if err != nil {
		logger.WithError(err).Warning("failed to verify signature")
		return nil, rpc.NewRpcSimpleValidationError(ErrInvalidSignature)
	}

	err = s.ds.Accounts.UpdatePublicKey(ctx, account, pkHex)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	token, err := CreateAuthToken(ctx, s.authSecret, account)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	err = s.ds.Accounts.RegenerateNonce(ctx, account)
	if err != nil {
		return nil, rpc.NewRpcInternalError(err)
	}

	resp := &v1.AuthResponse{Token: token}

	return resp, nil
}

func (s *Server) GetAccount(ctx context.Context, _ *v1.AccountRequest) (*v1.AccountResponse, error) {
	account, _ := rpcauth.AccountFromContext(ctx)
	resp := toAccountResponse(account)
	return resp, nil
}

func (s *Server) UpdateAccount(ctx context.Context, req *v1.UpdateAccountRequest) (*v1.AccountResponse, error) {
	account, _ := rpcauth.AccountFromContext(ctx)

	updateFields := new(datastore.UpdateAccountFields)

	if req.Username != nil &&
		len(req.Username.Value) > 0 &&
		account.Username.String != req.Username.Value {
		_, err := s.ds.Accounts.GetByUsername(ctx, req.Username.Value)
		if err == datastore.ErrAccountNotFound {
			updateFields.Username = pointer.ToString(req.Username.Value)
		} else if err != nil {
			return nil, rpc.NewRpcInternalError(err)
		} else {
			return nil, rpc.NewRpcSimpleValidationError(errors.New("username already used"))
		}
	}

	if req.Name != nil && len(req.Name.Value) > 0 {
		updateFields.Name = pointer.ToString(req.Name.Value)
	}

	if !updateFields.IsEmpty() {
		err := s.ds.Accounts.Update(ctx, account, *updateFields)
		if err != nil {
			return nil, rpc.NewRpcInternalError(err)
		}
	}

	resp := toAccountResponse(account)
	return resp, nil
}

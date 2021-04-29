package accounts

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/api/rpc"
	v1 "github.com/videocoin/marketplace/api/v1/accounts"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/rpcauth"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/pkg/grpcutil"
	"github.com/videocoin/marketplace/pkg/random"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
)

type Server struct {
	logger     *logrus.Entry
	validator  *grpcutil.RequestValidator
	authSecret string
	ds         *datastore.Datastore
	storage    *storage.Storage
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

	if req.ImageData != nil {
		link, err := s.handleImageData(req.ImageData.Value, account.ID)
		if err != nil {
			if err == ErrInvalidImageData {
				return nil, rpc.NewRpcSimpleValidationError(err)
			}
			return nil, rpc.NewRpcInternalError(err)
		}

		updateFields.ImageURL = pointer.ToString(link)
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

func (s *Server) handleImageData(data string, accountID int64) (string, error) {
	var (
		imageData    image.Image
		strImageData string
		isPng        bool
	)

	if strings.HasPrefix(data, "data:image/jpeg;base64,") {
		strImageData = strings.TrimPrefix(data, "data:image/jpeg;base64,")
	} else if strings.HasPrefix(data, "data:image/png;base64,") {
		strImageData = strings.TrimPrefix(data, "data:image/png;base64,")
		isPng = true
	} else {
		return "", ErrInvalidImageData
	}

	srcImageData, err := base64.StdEncoding.DecodeString(strImageData)
	if err != nil {
		return "", ErrInvalidImageData
	}

	imageReader := bytes.NewReader(srcImageData)
	if isPng {
		imageData, err = png.Decode(imageReader)
		if err != nil {
			return "", ErrInvalidImageData
		}
	} else {
		imageData, err = jpeg.Decode(imageReader)
		if err != nil {
			return "", ErrInvalidImageData
		}
	}

	resizedImage := resize.Resize(400, 0, imageData, resize.Lanczos2)
	croppedImage, err := cutter.Crop(resizedImage, cutter.Config{
		Width:  400,
		Height: 400,
		Mode: cutter.Centered,
	})

	rcImageJpeg := new(bytes.Buffer)
	err = jpeg.Encode(rcImageJpeg, croppedImage, nil)
	if err != nil {
		return "", err
	}

	imageID := random.RandomString(5)

	k := fmt.Sprintf("u/%d/r_%s.jpg", accountID, imageID)
	link, err := s.storage.PushPath(k, rcImageJpeg)
	if err != nil {
		return "", err
	}

	go func() {
		logger := s.logger.
			WithField("account_id", accountID).
			WithField("image_id", imageID)
		originalImageJpeg := new(bytes.Buffer)
		err = jpeg.Encode(originalImageJpeg, imageData, nil)
		if err != nil {
			logger.WithError(err).Error("failed to encode account image")
			return
		}

		k = fmt.Sprintf("u/%d/o_%s.jpg", accountID, imageID)
		_, err = s.storage.PushPath(k, originalImageJpeg)
		if err != nil {
			logger.WithError(err).Error("failed to push account image to storage")
			return
		}
	}()

	return link, nil
}

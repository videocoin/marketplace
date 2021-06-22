package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/AlekSi/pointer"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gocraft/dbr/v2"
	"github.com/goombaio/namegenerator"
	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"github.com/videocoin/marketplace/internal/auth"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/random"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"
	"time"
)

func (s *Server) getNonce(c echo.Context) error {
	address := strings.ToLower(c.Param("address"))
	if !ethcommon.IsHexAddress(address) {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": ErrInvalidAddress.Error()})
	}

	account, err := s.ds.Accounts.GetByAddress(context.Background(), address)
	if err != nil {
		if err == datastore.ErrAccountNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	resp := toNonceResponse(account)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) register(c echo.Context) error {
	req := new(RegisterRequest)
	err := c.Bind(req)
	if err != nil {
		return echo.ErrBadRequest
	}

	address := strings.ToLower(req.Address)
	if !ethcommon.IsHexAddress(address) {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": ErrInvalidAddress.Error()})
	}

	ctx := context.Background()
	_, err = s.ds.Accounts.GetByAddress(ctx, address)
	if err == nil {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": ErrAddressAlreadyRegistered.Error()})
	}
	if err != datastore.ErrAccountNotFound {
		return err
	}

	var username string

	for {
		seed := time.Now().UTC().UnixNano()
		nameGen := namegenerator.NewNameGenerator(seed)
		username = nameGen.Generate()

		_, err = s.ds.Accounts.GetByUsername(context.Background(), username)
		if err == datastore.ErrAccountNotFound {
			break
		}

		if err != nil {
			return err
		}

		time.Sleep(time.Millisecond*200)
	}

	account := &model.Account{
		Address:  address,
		Username: dbr.NewNullString(username),
	}
	err = s.ds.Accounts.Create(ctx, account)
	if err != nil {
		return err
	}

	resp := toRegisterResponse(account)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) auth(c echo.Context) error {
	req := new(AuthRequest)
	err := c.Bind(req)
	if err != nil {
		return echo.ErrBadRequest
	}

	address := strings.ToLower(req.Address)
	if !ethcommon.IsHexAddress(address) {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": ErrInvalidAddress.Error()})
	}

	logger := s.logger.
		WithField("address", address).
		WithField("signature", req.Signature)

	ctx := context.Background()
	account, err := s.ds.Accounts.GetByAddress(ctx, address)
	if err != nil {
		if err == datastore.ErrAccountNotFound {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": ErrAddressNotRegistered.Error()})
		}
		return err
	}

	pkHex, err := verifySignature(req.Signature, account.Nonce.String, account.Address)
	if err != nil {
		logger.WithError(err).Warning("failed to verify signature")
		return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": ErrInvalidSignature.Error()})
	}

	err = s.ds.Accounts.UpdatePublicKey(ctx, account, pkHex)
	if err != nil {
		return err
	}

	token, err := auth.CreateAuthToken(ctx, s.authSecret, account)
	if err != nil {
		return err
	}

	err = s.ds.Accounts.RegenerateNonce(ctx, account)
	if err != nil {
		return err
	}

	resp := &AuthResponse{Token: token}
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getAccount(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)
	resp := toAccountResponse(account)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) updateAccount(c echo.Context) error {
	ctxAccount := c.Get("account")
	account := ctxAccount.(*model.Account)

	req := new(UpdateAccountRequest)
	err := c.Bind(req)
	if err != nil {
		return echo.ErrBadRequest
	}

	updateFields := new(datastore.UpdateAccountFields)

	if req.Username != nil &&
		len(*req.Username) > 0 &&
		account.Username.String != *req.Username {
		_, err := s.ds.Accounts.GetByUsername(context.Background(), *req.Username)
		if err == datastore.ErrAccountNotFound {
			updateFields.Username = pointer.ToString(*req.Username)
		} else if err != nil {
			return err
		} else {
			return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": ErrUsernameAlreadyUsed.Error()})
		}
	}

	if req.Name != nil {
		updateFields.Name = pointer.ToString(*req.Name)
	}

	if req.Bio != nil {
		updateFields.Bio = pointer.ToString(*req.Bio)
	}

	if req.CustomURL != nil {
		updateFields.CustomURL = pointer.ToString(*req.CustomURL)
	}

	if req.YTUsername != nil {
		updateFields.YTUsername = pointer.ToString(*req.YTUsername)
	}

	if req.ImageData != nil {
		link, err := s.handleImageData(*req.ImageData, account.ID)
		if err != nil {
			if err == ErrInvalidImageData {
				return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": err.Error()})
			}
			return err
		}

		updateFields.ImageURL = pointer.ToString(link)
	}

	if req.CoverData != nil {
		link, err := s.handleCoverData(*req.CoverData, account.ID)
		if err != nil {
			if err == ErrInvalidImageData {
				return c.JSON(http.StatusPreconditionFailed, echo.Map{"message": err.Error()})
			}
			return err
		}

		updateFields.CoverURL = pointer.ToString(link)
	}

	if !updateFields.IsEmpty() {
		err := s.ds.Accounts.Update(context.Background(), account, *updateFields)
		if err != nil {
			return err
		}
	}

	resp := toAccountResponse(account)
	return c.JSON(http.StatusOK, resp)
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
		Mode:   cutter.Centered,
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

func (s *Server) handleCoverData(data string, accountID int64) (string, error) {
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

	resizedImage := resize.Resize(1680, 0, imageData, resize.Lanczos2)
	croppedImage, err := cutter.Crop(resizedImage, cutter.Config{
		Width:  1680,
		Height: 340,
		Mode:   cutter.Centered,
	})

	rcImageJpeg := new(bytes.Buffer)
	err = jpeg.Encode(rcImageJpeg, croppedImage, nil)
	if err != nil {
		return "", err
	}

	imageID := random.RandomString(5)

	k := fmt.Sprintf("u/%d/r_cover_%s.jpg", accountID, imageID)
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
			logger.WithError(err).Error("failed to encode account cvoer image")
			return
		}

		k = fmt.Sprintf("u/%d/o_cover_%s.jpg", accountID, imageID)
		_, err = s.storage.PushPath(k, originalImageJpeg)
		if err != nil {
			logger.WithError(err).Error("failed to push account cover image to storage")
			return
		}
	}()

	return link, nil
}

package api

import (
	"context"
	"github.com/AlekSi/pointer"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/wyvern"
	"net/http"
	"strconv"
)

func (s *Server) postOrder(c echo.Context) error {
	ctx := context.Background()
	req := new(PostOrderRequest)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	order := new(model.Order)
	order.WyvernOrder = new(wyvern.Order)
	err = copier.Copy(order.WyvernOrder, req)
	if err != nil {
		return err
	}

	if order.WyvernOrder == nil {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "invalid request")
	}

	if order.WyvernOrder.Metadata == nil || order.WyvernOrder.Metadata.Asset == nil {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "missing asset")
	}

	s.logger.Debugf("%+v\n", order.WyvernOrder)
	s.logger.Debugf("%+v\n", order.WyvernOrder.Maker)
	s.logger.Debugf("%+v\n", order.WyvernOrder.Taker)
	s.logger.Debugf("%+v\n", order.WyvernOrder.Metadata)
	s.logger.Debugf("%+v\n", order.WyvernOrder.Metadata.Asset)

	if order.WyvernOrder.Maker == nil || !ethcommon.IsHexAddress(order.WyvernOrder.Maker.Address) {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "invalid maker address")
	}

	if order.WyvernOrder.Taker == nil || !ethcommon.IsHexAddress(order.WyvernOrder.Taker.Address) {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "invalid taker address")
	}

	tokenID, _ := strconv.ParseInt(order.WyvernOrder.Metadata.Asset.ID, 10, 64)
	asset, err := s.ds.Assets.GetByTokenID(ctx, tokenID)
	if err != nil {
		if err == datastore.ErrAssetNotFound {
			return echo.NewHTTPError(http.StatusPreconditionFailed, "asset not found")
		}
		return err
	}

	maker, err := s.ds.Accounts.GetByAddress(ctx, order.WyvernOrder.Maker.Address)
	if err != nil {
		if err == datastore.ErrAccountNotFound {
			return echo.NewHTTPError(http.StatusPreconditionFailed, "maker not found")
		}
		return err
	}

	taker, err := s.ds.Accounts.GetByAddress(ctx, order.WyvernOrder.Taker.Address)
	if err != nil {
		if err == datastore.ErrAccountNotFound {
			return echo.NewHTTPError(http.StatusPreconditionFailed, "taker not found")
		}
		return err
	}

	feeRecipient, err := s.ds.Accounts.GetByAddress(ctx, order.WyvernOrder.FeeRecipient.Address)
	if err != nil {
		if err == datastore.ErrAccountNotFound {
			return echo.NewHTTPError(http.StatusPreconditionFailed, "fee recipient not found")
		}
		return err
	}

	order.MakerID = maker.ID
	order.TakerID = taker.ID
	order.OwnerID = asset.CreatedByID
	_ = copier.Copy(order.WyvernOrder.Maker, maker)
	_ = copier.Copy(order.WyvernOrder.Taker, taker)
	_ = copier.Copy(order.WyvernOrder.FeeRecipient, feeRecipient)

	err = s.ds.Orders.Create(ctx, order)
	if err != nil {
		return err
	}

	resp := new(OrderResponse)
	err = copier.Copy(resp, order.WyvernOrder)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getOrders(c echo.Context) error {
	ctx := context.Background()
	fltr := &datastore.OrderFilter{
		Sort: &datastore.SortOption{
			Field: "created_date",
			IsAsc: false,
		},
	}

	offset, _ := strconv.ParseUint(c.FormValue("offset"), 10, 64)
	limit, _ := strconv.ParseUint(c.FormValue("limit"), 10, 64)
	limitOpts := datastore.NewLimitOpts(offset, limit)

	orderBy := c.FormValue("order_by")
	orderDir := c.FormValue("order_direction")
	if orderBy != "" {
		fltr.Sort.Field = "created_date"
	}
	if orderDir == "asc" {
		fltr.Sort.IsAsc = true
	}

	reqSide := c.FormValue("side")
	if reqSide != "" {
		side, _ := strconv.Atoi(reqSide)
		fltr.Side = pointer.ToInt(side)
	}

	reqSaleKind := c.FormValue("sale_kind")
	if reqSaleKind != "" {
		saleKind, _ := strconv.Atoi(reqSaleKind)
		fltr.SaleKind = pointer.ToInt(saleKind)
	}

	reqAssetContractAddress := c.FormValue("asset_contract_address")
	if reqAssetContractAddress != "" {
		fltr.AssetContractAddress = pointer.ToString(reqAssetContractAddress)
	}

	reqTokenID := c.FormValue("token_id")
	if reqTokenID != "" {
		tokenID, _ := strconv.ParseInt(reqTokenID, 10, 64)
		if tokenID != 0 {
			asset, _ := s.ds.Assets.GetByTokenID(ctx, tokenID)
			if asset == nil {
				return returnEmptyOrders(c)
			}

			fltr.TokenID = pointer.ToInt64(asset.ID)
		}
	}

	reqPaymentTokenAddress := c.FormValue("payment_token_address")
	if reqPaymentTokenAddress != "" {
		fltr.PaymentTokenAddress = pointer.ToString(reqPaymentTokenAddress)
	}

	reqMaker := c.FormValue("maker")
	if reqMaker != "" && ethcommon.IsHexAddress(reqMaker){
		maker, _ := s.ds.Accounts.GetByAddress(ctx, reqMaker)
		if maker == nil {
			return returnEmptyOrders(c)
		}

		fltr.MakerID = pointer.ToInt64(maker.ID)
	}

	reqTaker := c.FormValue("taker")
	if reqTaker != "" && ethcommon.IsHexAddress(reqTaker){
		taker, _ := s.ds.Accounts.GetByAddress(ctx, reqTaker)
		if taker == nil {
			return returnEmptyOrders(c)
		}

		fltr.TakerID = pointer.ToInt64(taker.ID)
	}

	orders, err := s.ds.Orders.List(ctx, fltr, limitOpts)
	if err != nil {
		return err
	}

	tc, _ := s.ds.Orders.Count(ctx, fltr)
	countResp := &ItemsCountResponse{
		TotalCount: tc,
		Offset:     *limitOpts.Offset,
		Limit:      *limitOpts.Limit,
	}

	resp := toOrdersResponse(orders, countResp)

	return c.JSON(http.StatusOK, resp)
}

func returnEmptyOrders(c echo.Context) error {
	orders := make([]*model.Order, 0)
	return c.JSON(http.StatusOK, toOrdersResponse(orders, nil))
}
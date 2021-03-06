package api

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/auth"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/mediaprocessor"
	"github.com/videocoin/marketplace/internal/minter"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/pkg/logger"
	"net/http"
)

type Server struct {
	logger     *logrus.Entry
	addr       string
	authSecret string
	ds         *datastore.Datastore
	storage    *storage.Storage
	mp         *mediaprocessor.MediaProcessor
	e          *echo.Echo
	minter     *minter.Minter
}

func NewServer(ctx context.Context, opts ...ServerOption) (*Server, error) {
	srv := &Server{
		logger: ctxlogrus.Extract(ctx).WithField("system", "api"),
	}
	for _, o := range opts {
		if err := o(srv); err != nil {
			return nil, err
		}
	}

	srv.e = echo.New()
	srv.e.HideBanner = true
	srv.e.HidePort = true

	return srv, nil
}

func (s *Server) route() {
	logger.EchoLogger = s.logger

	s.e.Pre(middleware.RemoveTrailingSlash())
	s.e.Use(middleware.CORS())
	s.e.Use(logger.NewEchoLogrus())

	s.e.GET("/healthz", s.health)
	s.e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	v1 := s.e.Group("/api/v1")
	v1.GET("/swagger.json", s.getSwagger)

	registerGroup := v1.Group("/accounts")
	registerGroup.POST("", s.register)
	registerGroup.GET("/:address/nonce", s.getNonce)

	accountGroup := v1.Group("/account")
	accountGroup.Use(auth.JWTAuth(s.logger, s.ds, s.authSecret))
	accountGroup.GET("", s.getAccount)
	accountGroup.PUT("", s.updateAccount)

	authGroup := v1.Group("/auth")
	authGroup.POST("", s.auth)

	assetsGroup := v1.Group("/assets")
	assetsGroup.GET("", s.getAssets)
	assetsGroup.POST("", s.createAsset, auth.JWTAuth(s.logger, s.ds, s.authSecret))
	assetsGroup.GET("/:asset_id", s.getAsset)

	mediaGroup := v1.Group("/media")
	mediaGroup.POST("/upload", s.uploadMedia, auth.JWTAuth(s.logger, s.ds, s.authSecret))
	mediaGroup.GET("/:media_id", s.getMedia, auth.JWTAuth(s.logger, s.ds, s.authSecret))

	v1.GET("/asset/:contract_address/:token_id", s.getAssetByContractAddressAndTokenID)
	v1.GET("/tokens", s.getTokens)
	s.e.POST("/wyvern/v1/orders/post", s.postOrder, auth.JWTAuth(s.logger, s.ds, s.authSecret))
	s.e.GET("/wyvern/v1/orders", s.getOrders)

	myGroup := v1.Group("/my/assets")
	myGroup.Use(auth.JWTAuth(s.logger, s.ds, s.authSecret))
	myGroup.GET("", s.getMyAssets)
	myGroup.GET("/sold", s.getMySoldAssets)

	creatorsGroup := v1.Group("/creators")
	creatorsGroup.GET("", s.GetCreators)
	creatorsGroup.GET("/:creator_id", s.GetCreator)
	creatorsGroup.GET("/:creator_id/assets", s.getAssetsByCreator)

	spotlightGroup := v1.Group("/spotlight")
	spotlightGroup.GET("/assets/featured", s.getSpotlightFeaturedAssets)
	spotlightGroup.GET("/assets/live", s.getSpotlightLiveAssets)
	spotlightGroup.GET("/creators/featured", s.getSpotlightFeaturedCreators)

	activityGroup := v1.Group("/activity")
	activityGroup.Use(auth.JWTAuth(s.logger, s.ds, s.authSecret))
	activityGroup.GET("", s.getActivity)
}

func (s *Server) health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"alive": "OK"})
}

func (s *Server) Start(errCh chan error) {
	s.logger.WithField("addr", s.addr).Info("starting api server")

	s.route()

	go func() {
		err := s.e.Start(s.addr)
		if err == http.ErrServerClosed {
			err = nil
		}
		errCh <- err
	}()
}

func (s *Server) Stop() error {
	return s.e.Shutdown(context.Background())
}

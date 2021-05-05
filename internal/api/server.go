package api

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/kkdai/youtube/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echologrus "github.com/plutov/echo-logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/auth"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/mediaconverter"
	"github.com/videocoin/marketplace/internal/storage"
	"net/http"
)

type Server struct {
	logger     *logrus.Entry
	addr       string
	authSecret string
	gcpBucket  string
	ds         *datastore.Datastore
	storage    *storage.Storage
	mc         *mediaconverter.MediaConverter
	yt         *youtube.Client
	e          *echo.Echo
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

	srv.yt = &youtube.Client{}

	srv.e = echo.New()
	srv.e.HideBanner = true
	srv.e.HidePort = true

	return srv, nil
}

func (s *Server) route() {
	echologrus.Logger = s.logger.Logger

	s.e.Use(middleware.CORS())
	s.e.Use(echologrus.Hook())

	s.e.GET("/healthz", s.health)
	s.e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	v1 := s.e.Group("/api/v1")

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
	assetsGroup.Use(auth.JWTAuth(s.logger, s.ds, s.authSecret))
	assetsGroup.POST("/upload", s.upload)
	assetsGroup.POST("/ytupload", s.ytUpload)

	creatorsGroup := v1.Group("/creators")
	creatorsGroup.GET("", s.GetCreators)
	creatorsGroup.GET("/:creator_id", s.GetCreator)
	creatorsGroup.GET("/:creator_id/arts", s.getArtsByCreator)

	artsGroup := v1.Group("/arts")
	artsGroup.POST("", s.createArt, auth.JWTAuth(s.logger, s.ds, s.authSecret))
	artsGroup.GET("", s.getArts)
	artsGroup.GET("/:art_id", s.getArt)

	myGroup := v1.Group("/my/arts")
	myGroup.GET("", s.getMyArts)

	spotlightGroup := v1.Group("/spotlight")
	spotlightGroup.GET("/arts/featured", s.getSpotlightFeaturedArts)
	spotlightGroup.GET("/arts/live", s.getSpotlightLiveArts)
	spotlightGroup.GET("/creators/featured", s.getSpotlightFeaturedCreators)
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
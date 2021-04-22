package gateway

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echologrus "github.com/plutov/echo-logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/marketplace/api/v1/accounts"
	marketplacev1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/assets"
	"github.com/videocoin/marketplace/pkg/grpcutil"
	"github.com/videocoin/marketplace/pkg/jsonpb"
	"net/http"
)

type Gateway struct {
	logger      *logrus.Entry
	addr        string
	backendAddr string
	e           *echo.Echo
	gw          *runtime.ServeMux
	assetsSvc   *assets.AssetsService
}

func NewGateway(ctx context.Context, opts ...Option) (*Gateway, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	gw := &Gateway{
		e:      e,
		logger: ctxlogrus.Extract(ctx).WithField("system", "gateway"),
	}

	for _, o := range opts {
		if err := o(gw); err != nil {
			return nil, err
		}
	}

	marshaler := &jsonpb.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	annotators := []annotator{injectHeadersIntoMetadata}
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaler),
		runtime.WithMetadata(chainGrpcAnnotators(annotators...)),
		WithProtoHTTPErrorHandler(),
		runtime.WithForwardResponseOption(responseHeaderMatcher),
	)

	grpcOpts := grpcutil.DefaultClientDialOpts(gw.logger)

	err := accountsv1.RegisterAccountsServiceHandlerFromEndpoint(ctx, mux, gw.backendAddr, grpcOpts)
	if err != nil {
		return nil, err
	}

	err = marketplacev1.RegisterMarketplaceServiceHandlerFromEndpoint(ctx, mux, gw.backendAddr, grpcOpts)
	if err != nil {
		return nil, err
	}

	gw.gw = mux

	gw.route()

	return gw, nil
}

func (gw *Gateway) route() {
	echologrus.Logger = gw.logger.Logger

	gw.e.Use(middleware.CORS())
	gw.e.Use(echologrus.Hook())

	gw.e.GET("/healthz", gw.health)
	gw.e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	gw.e.Any("/api/v1/*", echo.WrapHandler(gw.gw))
	if gw.assetsSvc != nil {
		gw.assetsSvc.InitRoutes(gw.e)
	}
}

func (gw *Gateway) health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"alive": "OK"})
}

func (gw *Gateway) Start(errCh chan error) {
	gw.logger.WithField("addr", gw.addr).Info("starting gateway")

	go func() {
		err := gw.e.Start(gw.addr)
		if err == http.ErrServerClosed {
			err = nil
		}
		errCh <- err
	}()
}

func (gw *Gateway) Stop() error {
	return gw.e.Shutdown(context.Background())
}

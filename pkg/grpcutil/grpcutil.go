package grpcutil

import (
	"context"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpctracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"time"
)

func DefaultServerOpts(logger *logrus.Entry) []grpc.ServerOption {
	// grpclogrus.ReplaceGrpcLogger(logger)

	tracerOpts := []grpctracing.Option{
		grpctracing.WithTracer(opentracing.GlobalTracer()),
		grpctracing.WithFilterFunc(func(ctx context.Context, fullMethodName string) bool {
			if fullMethodName == "/grpc.health.v1.Health/Check" {
				return false
			}
			return true
		}),
	}
	logrusOpts := []grpclogrus.Option{
		grpclogrus.WithDecider(func(methodFullName string, err error) bool {
			if methodFullName == "/grpc.health.v1.Health/Check" {
				return false
			}
			return true
		}),
	}

	return []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpctracing.UnaryServerInterceptor(tracerOpts...),
			grpcprometheus.UnaryServerInterceptor,
			grpclogrus.UnaryServerInterceptor(logger, logrusOpts...),
			grpcvalidator.UnaryServerInterceptor(),
			grpcauth.UnaryServerInterceptor(auth),
		)),
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcctxtags.StreamServerInterceptor(),
			grpctracing.StreamServerInterceptor(tracerOpts...),
			grpcprometheus.StreamServerInterceptor,
			grpclogrus.StreamServerInterceptor(logger),
			grpcvalidator.StreamServerInterceptor(),
			grpcauth.StreamServerInterceptor(auth),
		)),
	}
}

func DefaultClientDialOpts(logger *logrus.Entry) []grpc.DialOption {
	tracerOpts := grpctracing.WithTracer(opentracing.GlobalTracer())

	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc.UnaryClientInterceptor(grpcmiddleware.ChainUnaryClient(
				grpctracing.UnaryClientInterceptor(tracerOpts),
				grpclogrus.UnaryClientInterceptor(logger),
				grpcprometheus.UnaryClientInterceptor,
			)),
		),
		grpc.WithStreamInterceptor(
			grpc.StreamClientInterceptor(grpcmiddleware.ChainStreamClient(
				grpctracing.StreamClientInterceptor(tracerOpts),
				grpclogrus.StreamClientInterceptor(logger),
				grpcprometheus.StreamClientInterceptor,
			)),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 10,
			Timeout:             time.Second * 10,
			PermitWithoutStream: true,
		}),
	}
}

func nopAuth(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func auth(ctx context.Context) (context.Context, error) {
	return ctx, grpc.Errorf(codes.Unauthenticated, "Authentication failed")
}

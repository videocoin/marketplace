package rpc

import (
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRpcInternal         = grpc.Errorf(codes.Internal, "Oops. Something went wrong! Sorry. We've let our engineers know.")
	ErrRpcUnauthenticated  = grpc.Errorf(codes.Unauthenticated, "Authentication failed")
	ErrRpcPermissionDenied = grpc.Errorf(codes.PermissionDenied, "Permission denied")
	ErrRpcNotFound         = grpc.Errorf(codes.NotFound, "Not found")
	ErrRpcBadRequest       = grpc.Errorf(codes.InvalidArgument, "Bad request")
)

func NewRpcValidationError(verr proto.Message) error {
	s, _ := status.New(codes.InvalidArgument, "invalid argument").WithDetails(verr)
	return s.Err()
}

func NewRpcInternalError(err error) error {
	s := status.New(codes.Internal, err.Error())
	return s.Err()
}

func NewRpcPermissionError(err error) error {
	s := status.New(codes.PermissionDenied, err.Error())
	return s.Err()
}

func NewRpcSimpleValidationError(err error) error {
	s := status.New(codes.FailedPrecondition, err.Error())
	return s.Err()
}

func IsNotFoundError(err error) bool {
	if s, ok := status.FromError(err); ok {
		if s.Code() == codes.NotFound {
			return true
		}
	}

	return false
}

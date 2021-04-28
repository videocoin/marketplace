package storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/textileio/go-threads/core/thread"
	bucketsd "github.com/textileio/textile/v2/api/bucketsd/client"
	bucketspb "github.com/textileio/textile/v2/api/bucketsd/pb"
	"github.com/textileio/textile/v2/api/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"time"
)

type Storage struct {
	config  *TextileConfig
	authCtx context.Context
	cli     *bucketsd.Client
	root    *bucketspb.RootResponse
}

func NewStorage(opts ...Option) (*Storage, error) {
	s := new(Storage)
	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithPerRPCCredentials(common.Credentials{}),
	}
	cli, err := bucketsd.NewClient("api.hub.textile.io:443", dialOpts...)
	if err != nil {
		return nil, err
	}
	s.cli = cli

	authCtx, err := common.CreateAPISigContext(
		common.NewAPIKeyContext(context.Background(), s.config.AuthKey),
		time.Now().Add(time.Minute),
		s.config.AuthSecret,
	)
	if err != nil {
		return nil, err
	}

	tid, err := thread.Decode(s.config.ThreadID)
	if err != nil {
		return nil, fmt.Errorf("invalid textile thread id: %s", err)
	}
	authCtx = common.NewThreadIDContext(authCtx, tid)

	s.authCtx = authCtx

	root, err := s.cli.Root(s.authCtx, s.config.BucketRootKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket root: %s", err)
	}

	s.root = root

	return s, nil
}

func (s *Storage) PushPath(path string, src io.Reader) (string, error) {
	_, _, err := s.cli.PushPath(s.authCtx, s.root.Root.Key, path, src)
	if err != nil {
		return "", err
	}

	links, err := s.cli.Links(s.authCtx, s.root.Root.Key, path)
	if err != nil {
		return "", err
	}

	return links.Ipns, nil
}

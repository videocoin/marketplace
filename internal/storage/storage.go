package storage

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/textileio/go-threads/core/thread"
	bucketsd "github.com/textileio/textile/v2/api/bucketsd/client"
	bucketspb "github.com/textileio/textile/v2/api/bucketsd/pb"
	"github.com/textileio/textile/v2/api/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"os"
	"time"
)

const (
	Textile    = "textile"
	NftStorage = "nftstorage"
)

var (
	ErrUnknownStorageBackend = errors.New("unknown storage backend")
)

type Storage struct {
	backend          string
	textileConfig    *TextileConfig
	nftStorageConfig *NftStorageConfig
	authCtx          context.Context
	ttCli            *bucketsd.Client
	ttRoot           *bucketspb.RootResponse
	nsCli            *NftStorageClient
}

func NewStorage(opts ...Option) (*Storage, error) {
	s := new(Storage)
	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}

	if s.nftStorageConfig != nil {
		s.backend = NftStorage
		s.nsCli = NewNftStorageClient(s.nftStorageConfig.ApiKey)
	}

	if s.textileConfig != nil {
		s.backend = Textile

		dialOpts := []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
			grpc.WithPerRPCCredentials(common.Credentials{}),
		}
		cli, err := bucketsd.NewClient("api.hub.textile.io:443", dialOpts...)
		if err != nil {
			return nil, err
		}
		s.ttCli = cli

		authCtx, err := common.CreateAPISigContext(
			common.NewAPIKeyContext(context.Background(), s.textileConfig.AuthKey),
			time.Now().Add(time.Minute),
			s.textileConfig.AuthSecret,
		)
		if err != nil {
			return nil, err
		}

		tid, err := thread.Decode(s.textileConfig.ThreadID)
		if err != nil {
			return nil, fmt.Errorf("invalid textile thread id: %s", err)
		}
		authCtx = common.NewThreadIDContext(authCtx, tid)

		s.authCtx = authCtx

		root, err := s.ttCli.Root(s.authCtx, s.textileConfig.BucketRootKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get bucket root: %s", err)
		}

		s.ttRoot = root
	}

	return s, nil
}

func (s *Storage) RootPath() string {
	if s.textileConfig != nil {
		return s.textileConfig.BucketRootKey
	}
	return ""
}

func (s *Storage) PushPath(path string, src io.Reader) (string, error) {
	if s.backend == NftStorage {
		return s.nsCli.PushPath(path, src)
	}

	if s.backend == Textile {
		result, _, err := s.ttCli.PushPath(s.authCtx, s.ttRoot.Root.Key, path, src)
		if err != nil {
			return "", err
		}

		return result.Cid().String(), nil
	}

	return "", ErrUnknownStorageBackend
}

func (s *Storage) Upload(input string, to string) (string, error) {
	f, err := os.Open(input)
	if err != nil {
		return "", err
	}

	cid, err := s.PushPath(to, f)
	if err != nil {
		return "", err
	}

	return cid, nil
}

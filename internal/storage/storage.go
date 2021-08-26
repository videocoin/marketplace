package storage

import (
	"bytes"
	gcpstorage "cloud.google.com/go/storage"
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
	gcpBucket        string
	authCtx          context.Context
	ttCli            *bucketsd.Client
	ttRoot           *bucketspb.RootResponse
	nsCli            *NftStorageClient
	gcpStorage       *gcpstorage.Client
	gcpBh            *gcpstorage.BucketHandle
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

	if s.gcpBucket != "" {
		gcpCli, err := gcpstorage.NewClient(context.Background())
		if err != nil {
			return nil, err
		}
		s.gcpStorage = gcpCli
		s.gcpBh = s.gcpStorage.Bucket(s.gcpBucket)
	}

	return s, nil
}

func (s *Storage) RootPath() string {
	if s.textileConfig != nil {
		return s.textileConfig.BucketRootKey
	}
	return ""
}

func (s *Storage) CacheRootPath() string {
	return s.gcpBucket
}

func (s *Storage) PushPath(path string, src io.Reader, public bool) (string, error) {
	var (
		buf bytes.Buffer
		cid string
	)

	r := io.TeeReader(src, &buf)

	if s.backend == NftStorage {
		result, err := s.nsCli.PushPath(path, r)
		if err != nil {
			return "", err
		}
		cid = result
	} else if s.backend == Textile {
		result, _, err := s.ttCli.PushPath(s.authCtx, s.ttRoot.Root.Key, path, r)
		if err != nil {
			return "", err
		}
		cid = result.Cid().String()
	} else {
		return "", ErrUnknownStorageBackend
	}

	if s.gcpBh != nil {
		w := s.gcpBh.Object(path).NewWriter(context.Background())
		if public {
			w.ACL = []gcpstorage.ACLRule{
				{
					Entity: gcpstorage.AllUsers,
					Role:   gcpstorage.RoleReader,
				},
			}
		}
		_, err := io.Copy(w, &buf)
		if err != nil {
			return "", err
		}
		_ = w.Close()
	}

	return cid, nil
}

func (s *Storage) Upload(input string, to string) (string, error) {
	f, err := os.Open(input)
	if err != nil {
		return "", err
	}
	defer f.Close()

	cid, err := s.PushPath(to, f, true)
	if err != nil {
		return "", err
	}

	return cid, nil
}

func (s *Storage) MultiUpload(inputs []string, to []string) (string, error) {
	cids := make([]string, 0)
	if len(inputs) != len(to) {
		return "", errors.New("different number of input/output paths")
	}

	if s.backend == Textile {
		for idx, input := range inputs {
			cid, err := s.Upload(input, to[idx])
			if err != nil {
				return "", err
			}
			cids = append(cids, cid)
		}
	}

	if s.backend == NftStorage {
		srcs := make([]io.Reader, 0)
		for _, input := range inputs {
			f, err := os.Open(input)
			if err != nil {
				return "", err
			}
			srcs = append(srcs, f)
		}
		cid, err := s.nsCli.PushPaths(inputs, srcs)
		if err != nil {
			return "", err
		}
		cids = append(cids, cid)
	}

	return cids[0], nil
}

func (s *Storage) MakePublic(path string) error {
	if s.gcpBh != nil {
		acl := s.gcpBh.Object(path).ACL()
		return acl.Set(context.Background(), gcpstorage.AllUsers, gcpstorage.RoleReader)
	}

	return nil
}

func (s *Storage) UploadToCloud(src io.Reader, path string) error {
	w := s.gcpBh.Object(path).NewWriter(context.Background())
	w.ACL = []gcpstorage.ACLRule{
		{
			Entity: gcpstorage.AllUsers,
			Role:   gcpstorage.RoleReader,
		},
	}
	_, err := io.Copy(w, src)
	if err != nil {
		return err
	}

	return w.Close()
}

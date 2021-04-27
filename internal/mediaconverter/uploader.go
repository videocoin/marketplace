package mediaconverter

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
)

func (mc *MediaConverter) uploadToStorage(ctx context.Context, src io.Reader, ct string, key string) error {
	w := mc.bh.Object(key).NewWriter(ctx)
	w.ContentType = ct
	w.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	if _, err := io.Copy(w, src); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

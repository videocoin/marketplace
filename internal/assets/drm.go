package assets

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	eciesgo "github.com/ecies/go"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/random"
	"io"
)

func GenerateDRMKeyID(account *model.Account) string {
	s := fmt.Sprintf("%d-%d-%s",
				account.ID,
				account.CreatedAt.UnixNano(),
				account.PublicKey.String)
	md5hash := md5.New()
	_, _ = io.WriteString(md5hash, s)
	return hex.EncodeToString(md5hash.Sum(nil))
}

func GenerateDRMKey(pkStr, ekStr string) (string, error) {
	pkBytes, err := hex.DecodeString(pkStr)
	if err != nil {
		return "", err
	}

	pk, err := eciesgo.NewPublicKeyFromBytes(pkBytes)
	if err != nil {
		return "", err
	}

	drmKey, err := eciesgo.Encrypt(pk, []byte(ekStr))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(drmKey), nil
}

func GenerateEncryptionKey() string {
	return hex.EncodeToString([]byte(random.RandomString(16)))
}

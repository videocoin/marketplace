package token

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kevinburke/nacl/randombytes"
	"github.com/twystd/tweetnacl-go/tweetnacl"
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

func GenerateDRMKey(pubKeyBase64, ekStr string) (string, error) {
	keyPair, err := tweetnacl.CryptoBoxKeyPair()
	if err != nil {
		return "", err
	}

	ephemPublicKey64 := base64.StdEncoding.EncodeToString(keyPair.PublicKey)

	pubKey, err := base64.StdEncoding.DecodeString(pubKeyBase64)
	if err != nil {
		return "", err
	}

	nonce := NewNonce()
	nonce64 := base64.StdEncoding.EncodeToString(nonce)

	cipher, err := tweetnacl.CryptoBox([]byte(ekStr), nonce, pubKey, keyPair.SecretKey)
	if err != nil {
		return "", err
	}

	cipher64 := base64.StdEncoding.EncodeToString(cipher)

	drmKey := map[string]string{
		"version": "x25519-xsalsa20-poly1305",
		"nonce": nonce64,
		"ephemPublicKey": ephemPublicKey64,
		"ciphertext": cipher64,
	}

	drmKeyJSON, err := json.Marshal(drmKey)
	if err != nil {
		return "", err
	}

	return common.Bytes2Hex(drmKeyJSON), nil
}

func GenerateEncryptionKey() string {
	return hex.EncodeToString([]byte(random.RandomString(16)))
}

func NewNonce() []byte {
	nonce := make([]byte, 24)
	randombytes.MustRead(nonce[:])
	return nonce
}

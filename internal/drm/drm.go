package drm

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kevinburke/nacl/randombytes"
	"github.com/twystd/tweetnacl-go/tweetnacl"
)

const drmXMLTemplate = `
<?xml version="1.0" encoding="UTF-8"?>
<GPACDRM type="CENC AES-CTR">
  <DRMInfo type="pssh" version="1">
    <BS ID128="1077efecc0b24d02ace33c1e52e2fb4b"/>
    <BS bits="32" value="1"/>
    <BS ID128="cd7eb9ff88f34caeb06185b00024e4c2"/>
  </DRMInfo>
  <CrypTrack IV_size="8" first_IV="%s" isEncrypted="1" saiSavedBox="senc">
    <key KID="%s" value="%s"/>
  </CrypTrack>
</GPACDRM>`

type Metadata struct {
	FirstIV string `json:"first_iv"`
	Key     string `json:"key"`
	KID     string `json:"kid"`
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewNonce() []byte {
	nonce := make([]byte, 24)
	randombytes.MustRead(nonce[:])
	return nonce
}

func GenerateFirstIV() string {
	s, _ := randomHex(8)
	return s
}

func GenerateKID() string {
	s, _ := randomHex(16)
	return s
}

func GenerateEncryptionKey() string {
	s, _ := randomHex(16)
	return s
}

func GenerateDRMKey(pubKeyBase64 string) (string, *Metadata, error) {
	meta := &Metadata{
		FirstIV: GenerateFirstIV(),
		Key:     GenerateEncryptionKey(),
		KID:     GenerateKID(),
	}

	message, err := json.Marshal(meta)
	if err != nil {
		return "", nil, err
	}

	keyPair, err := tweetnacl.CryptoBoxKeyPair()
	if err != nil {
		return "", nil, err
	}

	ephemPublicKey64 := base64.StdEncoding.EncodeToString(keyPair.PublicKey)

	pubKey, err := base64.StdEncoding.DecodeString(pubKeyBase64)
	if err != nil {
		return "", nil, err
	}

	nonce := NewNonce()
	nonce64 := base64.StdEncoding.EncodeToString(nonce)

	cipher, err := tweetnacl.CryptoBox(message, nonce, pubKey, keyPair.SecretKey)
	if err != nil {
		return "", nil, err
	}

	cipher64 := base64.StdEncoding.EncodeToString(cipher)

	drmKey := map[string]string{
		"version":        "x25519-xsalsa20-poly1305",
		"nonce":          nonce64,
		"ephemPublicKey": ephemPublicKey64,
		"ciphertext":     cipher64,
	}

	drmKeyJSON, err := json.Marshal(drmKey)
	if err != nil {
		return "", nil, err
	}

	return common.Bytes2Hex(drmKeyJSON), meta, nil
}

func GenerateDrmXml(meta *Metadata) string {
	return fmt.Sprintf(drmXMLTemplate, meta.FirstIV, meta.KID, meta.Key)
}

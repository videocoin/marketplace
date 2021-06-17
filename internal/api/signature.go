package api

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/videocoin/marketplace/pkg/ethutil"
	"strings"
)

func verifySignature(signature string, nonce string, address string) (string, error) {
	decodedSig, err := hexutil.Decode(signature)
	if err != nil {
		return "", err
	}

	if len(decodedSig) < 65 || (decodedSig[64] != 27 && decodedSig[64] != 28) {
		return "", ErrInvalidSignature
	}

	decodedSig[64] -= 27

	recoveredPublicKey, err := crypto.SigToPub(ethutil.SignHash([]byte(nonce)), decodedSig)
	if err != nil {
		return "", err
	}

	pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(recoveredPublicKey))

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPublicKey)

	verified := strings.ToLower(address) == strings.ToLower(recoveredAddress.String())
	if !verified {
		return "", ErrInvalidSignature
	}

	return pubKeyHex, nil
}

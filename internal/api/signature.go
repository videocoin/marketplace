package api

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/videocoin/marketplace/pkg/ethutil"
	"strings"
)

const (
	NoncePrefix = "Hi there from VideoCoin NFT! Sign this message to prove you have access to this wallet and we’ll log you in. This won’t cost you any coins.\nTo stop hackers using your wallet, here’s a unique message ID they can’t guess: "
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

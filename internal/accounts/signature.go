package accounts

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
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

	recoveredPublicKey, err := crypto.SigToPub(signHash([]byte(nonce)), decodedSig)
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

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

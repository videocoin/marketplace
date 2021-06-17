package ethutil

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func SignHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}


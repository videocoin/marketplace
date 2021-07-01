package ethutil

import (
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

func ParseBigInt(value string) (*big.Int, error) {
	f := new(big.Int)
	_, err := fmt.Sscan(value, f)
	return f, err
}

func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236)  //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236)  //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

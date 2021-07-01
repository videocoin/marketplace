package main

import (
	"encoding/hex"
	"fmt"
	eciesgo "github.com/ecies/go/v2"
	"os"
)

func main() {
	privk, err := eciesgo.NewPrivateKeyFromHex(os.Getenv("PRIV_KEY"))
	if err != nil {
		panic(err)
	}

	drmKeyBytes, err := hex.DecodeString(os.Getenv("DRM_KEY"))
	if err != nil {
		panic(err)
	}

	ek, err := eciesgo.Decrypt(privk, drmKeyBytes)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Encryption Key: %s\n", string(ek))
}

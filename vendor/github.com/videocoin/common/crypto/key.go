package crypto

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// DecryptKeyFile ...
func DecryptKeyFile(keyFilePath, passFilePath string) (*keystore.Key, error) {
	f, err := os.Open(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open keyfile %s: %v", keyFilePath, err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	passphrase := ""
	if len(passFilePath) != 0 {
		f, err = os.Open(passFilePath)
		if err != nil {
			return nil, fmt.Errorf("can't open %s: %v", passFilePath, err)
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}

		passphrase = strings.TrimRightFunc(string(data), func(r rune) bool {
			return r == '\r' || r == '\n'
		})
	}

	key, err := keystore.DecryptKey(data, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt a key %s: %v", keyFilePath, err)
	}

	return key, nil
}

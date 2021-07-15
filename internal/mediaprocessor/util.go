package mediaprocessor

import (
	"encoding/hex"
	"math/rand"
	"os"
	"path/filepath"
)

func genTempFilepath(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

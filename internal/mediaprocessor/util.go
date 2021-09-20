package mediaprocessor

import (
	"cloud.google.com/go/storage"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

func genTempFilepath(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

func downloadFile(url string, filepath string, ir *storage.Reader) (err error) {
	out, err := os.Create(filepath)
	if err != nil  {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	var src io.Reader
	if ir == nil {
		src = resp.Body
	} else {
		src = ir
	}

	_, err = io.Copy(out, src)
	if err != nil  {
		return err
	}

	return nil
}

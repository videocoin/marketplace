package random

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomStringWithCharset(length int, charset string) string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return RandomStringWithCharset(length, charset)
}

func StringWithTS(length int) string {
	return fmt.Sprintf(
		"%s-%s",
		RandomString(6),
		strconv.FormatInt(time.Now().UnixNano(), 10),
	)
}

package apiServer

import (
	"crypto/md5"
	"encoding/hex"
)

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

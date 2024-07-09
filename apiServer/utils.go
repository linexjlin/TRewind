package apiServer

import (
	"crypto/md5"
	"encoding/hex"
)

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func extractFilenameAndExtra(input string) (name, extra string) {
	for i, r := range input {
		if r == '\n' && len(name) != 0 {
			// left input is left
			extra = input[i+1:]
			break
		}
		name += string(r)
	}
	return name, extra
}

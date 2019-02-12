package server

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func getMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// resized file doesnt exist
		return false
	}
	return true
}

package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

func SaltPassword(password, salt string) string {
	var (
		hasher hash.Hash
	)

	hasher = sha256.New()
	hasher.Write([]byte(password))
	hasher.Write([]byte(salt))

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

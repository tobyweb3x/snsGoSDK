package utils

import (
	"crypto/sha256"
	spl "snsGoSDK/spl"
)

func GetHashedNameSync(name string) []byte {
	hashed := sha256.Sum256([]byte(spl.HashPrefix + name))
	return hashed[:]
}

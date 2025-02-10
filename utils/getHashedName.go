package utils

import (
	"crypto/sha256"
	spl "snsGoSDK/spl"
)

func GetHashedNameSync(name string) []byte {
	input := spl.HashPrefix + name
	hashed := sha256.Sum256([]byte(input))
	return hashed[:]
}

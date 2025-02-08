package spl_name_services

import "crypto/sha256"

func GetHashedNameSync(name string) []byte {
	input := HashPrefix + name
	hashed := sha256.Sum256([]byte(input))
	return hashed[:]
}

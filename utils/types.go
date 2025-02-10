package utils

import (
	"github.com/gagliardetto/solana-go"
)

// DomainKeyResult is a helper struct for GetDomainKeySync.
type DomainKeyResult struct {
	PubKey      solana.PublicKey
	Parent      solana.PublicKey
	Hashed      []byte
	IsSub       bool
	IsSubRecord bool
}

// deriveResult is a helper struct for GetDomainKeySync.
type deriveResult struct {
	PubKey solana.PublicKey
	Hashed []byte
}

// RecordVersion type for spl.

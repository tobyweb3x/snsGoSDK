package record

import (
	"crypto/ed25519"

	"github.com/gagliardetto/solana-go"
)

func CheckSolRecord(
	record, signedRecord []byte,
	pubkey solana.PublicKey,
) bool {
	return ed25519.Verify(pubkey.Bytes(), record, signedRecord)
}

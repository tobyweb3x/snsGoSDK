package utils_test

import (
	"crypto/rand"
	"crypto/sha256"
	"snsGoSDK/utils"
	"testing"

	"github.com/gagliardetto/solana-go"
)

func BenchmarkGetNameAccountKeySync(b *testing.B) {
	type input struct {
		hashedName []byte
		nameClass  solana.PublicKey
		nameParent solana.PublicKey
	}
	inputs := make([]input, 0, 30)

	for i := 0; i < 30; i++ {
		raw := make([]byte, 32)
		if _, err := rand.Read(raw); err != nil {
			b.Fatalf("failed to generate random input: %v", err)
		}
		hashed := sha256.Sum256(raw)

		var class, parent solana.PublicKey
		if i%3 != 0 {
			class = solana.NewWallet().PublicKey()
		}
		if i%5 != 0 {
			parent = solana.NewWallet().PublicKey()
		}

		inputs = append(inputs, input{
			hashedName: hashed[:],
			nameClass:  class,
			nameParent: parent,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := utils.GetNameAccountKeySync(inputs[i%30].hashedName, inputs[i%30].nameClass, inputs[i%30].nameParent)
		if err != nil {
			b.Fatalf("GetNameAccountKeySync failed: %v", err)
		}
	}
}

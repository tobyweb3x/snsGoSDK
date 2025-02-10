package utils

import (
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
)

func GetNameAccountKeySync(hashedName []byte, nameClass, nameParent solana.PublicKey) (solana.PublicKey, uint8, error) {
	var seeds [][]byte
	seeds = append(seeds, hashedName)

	if nameClass.IsZero() {
		seeds = append(seeds, make([]byte, 32))
	} else {
		seeds = append(seeds, nameClass.Bytes())
	}

	if nameParent.IsZero() {
		seeds = append(seeds, make([]byte, 32))
	} else {
		seeds = append(seeds, nameParent.Bytes())
	}

	nameAccountKey, nonce, err := solana.FindProgramAddress(seeds, spl.NameProgramID)

	if err != nil {
		return solana.PublicKey{}, nonce, err
	}

	return nameAccountKey, nonce, nil
}

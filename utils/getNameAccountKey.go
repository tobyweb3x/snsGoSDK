package utils

import (
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
)

func GetNameAccountKeySync(hashedName []byte, nameClass, nameParent solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := make([][]byte, 0, 3)
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

	return solana.FindProgramAddress(seeds, spl.NameProgramID)
}

func GetNameAccountKeySync2(hashedName []byte, nameClass, nameParent solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		hashedName,
		make([]byte, 32),
		make([]byte, 32),
	}

	if !nameClass.IsZero() {
		seeds[1] = nameClass.Bytes()
	}

	if !nameParent.IsZero() {
		seeds[2] = nameParent.Bytes()
	}

	return solana.FindProgramAddress(seeds, spl.NameProgramID)
}

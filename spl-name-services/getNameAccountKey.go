package spl_name_services

import (
	"github.com/gagliardetto/solana-go"
)

func GetNameAccountKeySync(hashedName []byte, nameClass, nameParent solana.PublicKey) (solana.PublicKey, uint8, error) {
	var seeds [][]byte
	seeds = append(seeds, hashedName)

	if nameClass == NoPublickKeyArg {
		seeds = append(seeds, make([]byte, 32))
	} else {
		seeds = append(seeds, nameClass.Bytes())
	}

	if nameParent == NoPublickKeyArg {
		seeds = append(seeds, make([]byte, 32))
	} else {
		seeds = append(seeds, nameParent.Bytes())
	}

	nameAccountKey, nonce, err := solana.FindProgramAddress(seeds, NameProgramID)

	if err != nil {
		return solana.PublicKey{}, nonce, err
	}

	return nameAccountKey, nonce, nil
}

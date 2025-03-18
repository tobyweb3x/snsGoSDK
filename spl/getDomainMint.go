package spl

import "github.com/gagliardetto/solana-go"

func getDomainMint(domain solana.PublicKey) (solana.PublicKey, error) {
	mint, _, err := solana.FindProgramAddress(
		[][]byte{MIntPrefix, domain.Bytes()},
		NameTokenizerID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return mint, nil
}

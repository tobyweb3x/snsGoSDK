package spl_name_services

import "github.com/gagliardetto/solana-go"

func getDomainMint(domain solana.PublicKey) (solana.PublicKey, error) {
	mint, _, err := solana.FindProgramAddress(
		[][]byte{MINT_PREFIX, domain.Bytes()},
		NAME_TOKENIZER_ID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return mint, nil
}

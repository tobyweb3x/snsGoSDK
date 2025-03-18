package nft

import (
	"snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
)

func GetDomainMint(domain solana.PublicKey) (solana.PublicKey, uint8, error) {
	return solana.FindProgramAddress(
		[][]byte{
			spl.MIntPrefix,
			domain.Bytes(),
		},
		spl.NameTokenizerID,
	)
}

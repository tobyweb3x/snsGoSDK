package bindings

import (
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
)

func RegisterWithNft(
	name string,
	space uint32,
	nameAccount,
	reverseLookupAccount,
	buyer,
	nftSource,
	nftMetadata,
	nftMint,
	masterEdition solana.PublicKey,
) (*solana.GenericInstruction, error) {

	state, _, err := solana.FindProgramAddress(
		[][]byte{nameAccount.Bytes()},
		spl.RegisterProgramID,
	)
	if err != nil {
		return nil, err
	}
	return instructions.NewCreateWithNFTInstruction(
		name,
		space,
	).GetInstruction(
		spl.RegisterProgramID,
		spl.NameProgramID,
		spl.RootDomainAccount,
		nameAccount,
		reverseLookupAccount,
		solana.SystemProgramID,
		spl.ReverseLookupClass,
		buyer,
		nftSource,
		nftMetadata,
		nftMint,
		masterEdition,
		spl.WolvesCollectionMetadata,
		solana.TokenProgramID,
		solana.SysVarRentPubkey,
		state,
		spl.MetaplexID,
	)
}

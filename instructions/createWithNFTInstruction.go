package instructions

import (
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type CreateWithNFTInstruction struct {
	Tag   uint8
	Name  string
	Space uint32
}

func (rf *CreateWithNFTInstruction) serialize() ([]byte, error) {
	return borsh.Serialize(*rf)
}

func NewCreateWithNFTInstruction(name string, space uint32) *CreateWithNFTInstruction {
	return &CreateWithNFTInstruction{
		Tag:   17,
		Name:  name,
		Space: space,
	}
}

func (rf *CreateWithNFTInstruction) getInstruction(
	programId,
	namingServiceProgram,
	rootDomain,
	name,
	reverseLookup,
	systemProgram,
	centralState,
	buyer,
	nftSource,
	nftMetadata,
	nftMint,
	masterEdition,
	collection,
	splTokenProgram,
	rentSysvar,
	state,
	mplTokenMetadata solana.PublicKey) *solana.GenericInstruction {

	data, err := rf.serialize()
	if err != nil {
		panic(err)
	}

	key := solana.AccountMetaSlice{
		{PublicKey: namingServiceProgram, IsSigner: false, IsWritable: false},
		{PublicKey: rootDomain, IsSigner: false, IsWritable: false},
		{PublicKey: name, IsSigner: false, IsWritable: true},
		{PublicKey: reverseLookup, IsSigner: false, IsWritable: true},
		{PublicKey: systemProgram, IsSigner: false, IsWritable: false},
		{PublicKey: centralState, IsSigner: false, IsWritable: false},
		{PublicKey: buyer, IsSigner: true, IsWritable: true},
		{PublicKey: nftSource, IsSigner: false, IsWritable: true},
		{PublicKey: nftMetadata, IsSigner: false, IsWritable: true},
		{PublicKey: nftMint, IsSigner: false, IsWritable: true},
		{PublicKey: masterEdition, IsSigner: false, IsWritable: true},
		{PublicKey: collection, IsSigner: false, IsWritable: true},
		{PublicKey: splTokenProgram, IsSigner: false, IsWritable: false},
		{PublicKey: rentSysvar, IsSigner: false, IsWritable: false},
		{PublicKey: state, IsSigner: false, IsWritable: false},
		{PublicKey: mplTokenMetadata, IsSigner: false, IsWritable: false},
	}

	return solana.NewInstruction(
		programId,
		key,
		data,
	)
}

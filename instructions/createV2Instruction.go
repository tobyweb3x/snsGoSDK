package instructions

import (
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type CreateV2Instruction struct {
	Tag   uint8
	Name  string
	Space uint32
}

func (c *CreateV2Instruction) Serialize() ([]byte, error) {
	return borsh.Serialize(*c)
}

func NewCreateV2Instruction(name string, space uint32) *CreateV2Instruction {
	return &CreateV2Instruction{
		Tag:   9,
		Name:  name,
		Space: space,
	}
}

func (c *CreateV2Instruction) GetInstruction(
	programId,
	rentSysvarAccount,
	nameProgramId,
	rootDomain,
	nameAccount,
	reverseLookupAccount,
	centralState,
	buyer,
	buyerTokenAccount,
	usdcVault,
	state solana.PublicKey) (*solana.GenericInstruction, error) {

	data, err := c.Serialize()
	if err != nil {
		return nil, err
	}

	keys := solana.AccountMetaSlice{
		{PublicKey: rentSysvarAccount, IsSigner: false, IsWritable: false},
		{PublicKey: nameProgramId, IsSigner: false, IsWritable: false},
		{PublicKey: rootDomain, IsSigner: false, IsWritable: false},
		{PublicKey: nameAccount, IsSigner: false, IsWritable: true},
		{PublicKey: reverseLookupAccount, IsSigner: false, IsWritable: true},
		{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		{PublicKey: centralState, IsSigner: false, IsWritable: false},
		{PublicKey: buyer, IsSigner: true, IsWritable: true},
		{PublicKey: buyerTokenAccount, IsSigner: false, IsWritable: true},
		{PublicKey: usdcVault, IsSigner: false, IsWritable: true},
		{PublicKey: solana.TokenProgramID, IsSigner: false, IsWritable: false},
		{PublicKey: state, IsSigner: false, IsWritable: false},
	}

	return solana.NewInstruction(
		programId,
		keys,
		data,
	), nil
}

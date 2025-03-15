package instructions

import (
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type CreateSplitV2Instruction struct {
	Tag            uint8
	Name           string
	Space          uint32
	ReferrerIdxOpt *uint16
}

func NewCreateSplitV2Instruction(name string, space uint32, referrerIdxOpt *uint16) *CreateSplitV2Instruction {
	return &CreateSplitV2Instruction{
		Tag:            20,
		Name:           name,
		Space:          space,
		ReferrerIdxOpt: referrerIdxOpt,
	}
}

func (cs *CreateSplitV2Instruction) GetInstruction(
	programId,
	namingServiceProgram,
	rootDomain,
	name,
	reverseLookup,
	systemProgram,
	centralState,
	buyer,
	domainOwner,
	feePayer,
	buyerTokenSource,
	pythFeedAccount,
	vault,
	splTokenProgram,
	rentSysvar,
	state solana.PublicKey,
	referrerAccountOpt solana.PublicKey,
) (*solana.GenericInstruction, error) {

	data, err := borsh.Serialize(*cs)
	if err != nil {
		return nil, err
	}

	keys := solana.AccountMetaSlice{
		{PublicKey: namingServiceProgram, IsSigner: false, IsWritable: false},
		{PublicKey: rootDomain, IsSigner: false, IsWritable: false},
		{PublicKey: name, IsSigner: false, IsWritable: true},
		{PublicKey: reverseLookup, IsSigner: false, IsWritable: true},
		{PublicKey: systemProgram, IsSigner: false, IsWritable: false},
		{PublicKey: centralState, IsSigner: false, IsWritable: false},
		{PublicKey: buyer, IsSigner: true, IsWritable: true},
		{PublicKey: domainOwner, IsSigner: false, IsWritable: false},
		{PublicKey: feePayer, IsSigner: true, IsWritable: true},
		{PublicKey: buyerTokenSource, IsSigner: false, IsWritable: true},
		{PublicKey: pythFeedAccount, IsSigner: false, IsWritable: false},
		{PublicKey: vault, IsSigner: false, IsWritable: true},
		{PublicKey: splTokenProgram, IsSigner: false, IsWritable: false},
		{PublicKey: rentSysvar, IsSigner: false, IsWritable: false},
		{PublicKey: state, IsSigner: false, IsWritable: false},
	}

	if !referrerAccountOpt.IsZero() {
		keys = append(keys, &solana.AccountMeta{PublicKey: referrerAccountOpt, IsSigner: false, IsWritable: true})
	}

	return solana.NewInstruction(
		programId,
		keys,
		data,
	), nil
}

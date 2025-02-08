package instructions

import (
	"bytes"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type CreateInstructionV3 struct {
	Tag            uint8
	Name           string
	Space          uint32
	ReferrerIdxOpt *uint16
}

func (civ3 *CreateInstructionV3) Serialize() ([]byte, error) {
	return borsh.Serialize(*civ3)
}

func NewCreateInstructionV3(name string, space uint32, referrerIdxOpt *uint16) *CreateInstructionV3 {
	return &CreateInstructionV3{
		Tag:            13,
		Name:           name,
		Space:          space,
		ReferrerIdxOpt: referrerIdxOpt,
	}
}

func (civ3 *CreateInstructionV3) GetInstruction(
	programId,
	namingServiceProgram,
	rootDomain,
	name,
	reverseLookup,
	systemProgram,
	centralState,
	buyer,
	buyerTokenSource,
	pythMappingAcc,
	pythProductAcc,
	pythPriceAcc,
	vault,
	splTokenProgram,
	rentSysvar,
	state solana.PublicKey,
	referrerAccountOpt *solana.PublicKey,
) *solana.GenericInstruction {

	data, err := civ3.Serialize()
	if err != nil {
		panic(err)
	}

	var dataBuffer bytes.Buffer
	dataBuffer.Write(data)

	keys := solana.AccountMetaSlice{
		{PublicKey: namingServiceProgram, IsSigner: false, IsWritable: false},
		{PublicKey: rootDomain, IsSigner: false, IsWritable: false},
		{PublicKey: name, IsSigner: false, IsWritable: true},
		{PublicKey: reverseLookup, IsSigner: false, IsWritable: true},
		{PublicKey: systemProgram, IsSigner: false, IsWritable: false},
		{PublicKey: centralState, IsSigner: false, IsWritable: false},
		{PublicKey: buyer, IsSigner: true, IsWritable: true},
		{PublicKey: buyerTokenSource, IsSigner: false, IsWritable: true},
		{PublicKey: pythMappingAcc, IsSigner: false, IsWritable: false},
		{PublicKey: pythProductAcc, IsSigner: false, IsWritable: false},
		{PublicKey: pythPriceAcc, IsSigner: false, IsWritable: false},
		{PublicKey: vault, IsSigner: false, IsWritable: true},
		{PublicKey: splTokenProgram, IsSigner: false, IsWritable: false},
		{PublicKey: rentSysvar, IsSigner: false, IsWritable: false},
		{PublicKey: state, IsSigner: false, IsWritable: false},
	}

	if !referrerAccountOpt.IsZero() {
		keys = append(keys, &solana.AccountMeta{PublicKey: *referrerAccountOpt, IsSigner: false, IsWritable: true})
	}

	return solana.NewInstruction(
		programId,
		keys,
		dataBuffer.Bytes(),
	)
}

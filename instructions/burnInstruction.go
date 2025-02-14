package instructions

import (
	"bytes"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type BurnInstruction struct {
	Tag uint8
}

func NewBurnInstruction() *BurnInstruction {
	return &BurnInstruction{
		Tag: 16,
	}
}

func (bi *BurnInstruction) GetInstruction(
	programId,
	nameServiceId,
	systemProgram,
	domain,
	reverse,
	resellingState,
	state,
	centralState,
	owner,
	target solana.PublicKey,
) (*solana.GenericInstruction, error) {

	data, err := borsh.Serialize(*bi)
	if err != nil {
		return nil, err
	}

	var dataBuffer bytes.Buffer
	dataBuffer.Write(data)

	keys := solana.AccountMetaSlice{
		{PublicKey: nameServiceId, IsSigner: false, IsWritable: false},
		{PublicKey: systemProgram, IsSigner: false, IsWritable: false},
		{PublicKey: domain, IsSigner: false, IsWritable: true},
		{PublicKey: reverse, IsSigner: false, IsWritable: true},
		{PublicKey: resellingState, IsSigner: false, IsWritable: true},
		{PublicKey: state, IsSigner: false, IsWritable: true},
		{PublicKey: centralState, IsSigner: false, IsWritable: false},
		{PublicKey: owner, IsSigner: true, IsWritable: false},
		{PublicKey: target, IsSigner: false, IsWritable: true},
	}

	return solana.NewInstruction(
		programId,
		keys,
		dataBuffer.Bytes(),
	), nil
}

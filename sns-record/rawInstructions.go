package snsRecord

import (
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type AllocateAndPostRecordInstruction struct {
	Tag     uint8
	Record  string
	Content []byte
}

func NewAllocateAndPostRecordInstruction(record string, content []byte) *AllocateAndPostRecordInstruction {
	return &AllocateAndPostRecordInstruction{
		Tag:     1,
		Record:  record,
		Content: content,
	}
}

func (apr *AllocateAndPostRecordInstruction) GetInstruction(
	programId,
	systemProgram,
	splNameServiceProgram,
	feePayer,
	record,
	domain,
	domainOwner,
	centralState solana.PublicKey,
) (*solana.GenericInstruction, error) {

	data, err := borsh.Serialize(*apr)
	if err != nil {
		return nil, err
	}

	keys := solana.AccountMetaSlice{
		{PublicKey: systemProgram, IsSigner: false, IsWritable: false},
		{PublicKey: splNameServiceProgram, IsSigner: false, IsWritable: false},
		{PublicKey: feePayer, IsSigner: true, IsWritable: true},
		{PublicKey: record, IsSigner: false, IsWritable: true},
		{PublicKey: domain, IsSigner: false, IsWritable: true},
		{PublicKey: domainOwner, IsSigner: true, IsWritable: true},
		{PublicKey: centralState, IsSigner: false, IsWritable: false},
	}

	return solana.NewInstruction(
		programId,
		keys,
		data,
	), nil
}

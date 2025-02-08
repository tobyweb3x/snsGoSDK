package instructions

import (
	"bytes"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type CreateReverseInstruction struct {
	Tag  uint8
	Name string
}

func (cri *CreateReverseInstruction) Serialize() ([]byte, error) {
	return borsh.Serialize(*cri)
}

func NewCreateReverseInstruction(name string) *CreateReverseInstruction {
	return &CreateReverseInstruction{
		Tag:  12,
		Name: name,
	}
}

func (cri *CreateReverseInstruction) GetInstruction(
	programId,
	namingServiceProgram,
	rootDomain,
	reverseLookup,
	systemProgram,
	centralState,
	feePayer,
	rentSysvar solana.PublicKey,
	parentName *solana.PublicKey,
	parentNameOwner *solana.PublicKey,
) *solana.GenericInstruction {

	data, err := cri.Serialize()
	if err != nil {
		panic(err)
	}

	var dataBuffer bytes.Buffer
	dataBuffer.Write(data)

	keys := solana.AccountMetaSlice{
		{PublicKey: namingServiceProgram, IsSigner: false, IsWritable: false},
		{PublicKey: rootDomain, IsSigner: false, IsWritable: false},
		{PublicKey: reverseLookup, IsSigner: false, IsWritable: true},
		{PublicKey: systemProgram, IsSigner: false, IsWritable: false},
		{PublicKey: centralState, IsSigner: false, IsWritable: false},
		{PublicKey: feePayer, IsSigner: true, IsWritable: true},
		{PublicKey: rentSysvar, IsSigner: false, IsWritable: false},
	}

	if parentName != nil {
		keys = append(keys, &solana.AccountMeta{PublicKey: *parentName, IsSigner: false, IsWritable: true})
	}

	if parentNameOwner != nil {
		keys = append(keys, &solana.AccountMeta{PublicKey: *parentNameOwner, IsSigner: true, IsWritable: true})
	}

	return solana.NewInstruction(
		programId,
		keys,
		dataBuffer.Bytes(),
	)
}

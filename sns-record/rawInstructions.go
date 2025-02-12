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

type DeleteRecordInstruction struct {
	Tag uint8
}

func NewDeleteRecordInstruction() *DeleteRecordInstruction {
	return &DeleteRecordInstruction{
		Tag: 5,
	}
}

func (dri *DeleteRecordInstruction) GetInstruction(
	programId,
	systemProgram,
	splNameServiceProgram,
	feePayer,
	record,
	domain,
	domainOwner,
	centralState solana.PublicKey,
) (*solana.GenericInstruction, error) {
	data, err := borsh.Serialize(*dri)
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

type ValidateEthereumSignatureInstruction struct {
	Tag            uint8
	Validation     uint8
	Signature      []byte
	ExpectedPubkey []byte
}

func NewValidateEthereumSignatureInstruction(validation uint8, signature []byte, expectedPubkey []byte) *ValidateEthereumSignatureInstruction {
	return &ValidateEthereumSignatureInstruction{
		Tag:            4,
		Validation:     validation,
		Signature:      signature,
		ExpectedPubkey: expectedPubkey,
	}
}

func (vesi *ValidateEthereumSignatureInstruction) GetInstruction(
	programId,
	systemProgram,
	splNameServiceProgram,
	feePayer,
	record,
	domain,
	domainOwner,
	centralState solana.PublicKey,
) (*solana.GenericInstruction, error) {
	data, err := borsh.Serialize(*vesi)
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

type EditRecordInstruction struct {
	Tag     uint8
	Record  string
	Content []byte
}

func NewEditRecordInstruction(record string, content []byte) *EditRecordInstruction {
	return &EditRecordInstruction{
		Tag:     2,
		Record:  record,
		Content: content,
	}
}

func (eri *EditRecordInstruction) GetInstruction(
	programId,
	systemProgram,
	splNameServiceProgram,
	feePayer,
	record,
	domain,
	domainOwner,
	centralState solana.PublicKey,
) (*solana.GenericInstruction, error) {

	data, err := borsh.Serialize(*eri)
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

type ValidateSolanaSignatureInstruction struct {
	Tag       uint8
	Staleness bool
}

func NewValidateSolanaSignatureInstruction(staleness bool) *ValidateSolanaSignatureInstruction {
	return &ValidateSolanaSignatureInstruction{
		Tag:       3,
		Staleness: staleness,
	}
}

func (vssi *ValidateSolanaSignatureInstruction) GetInstruction(
	programId,
	systemProgram,
	splNameServiceProgram,
	feePayer,
	record,
	domain,
	domainOwner,
	centralState,
	verifier solana.PublicKey,
) (*solana.GenericInstruction, error) {
	data, err := borsh.Serialize(*vssi)
	if err != nil {
		return nil, err
	}

	keys := solana.AccountMetaSlice{
		{PublicKey: systemProgram, IsSigner: false, IsWritable: false},
		{PublicKey: splNameServiceProgram, IsSigner: false, IsWritable: false},
		{PublicKey: feePayer, IsSigner: true, IsWritable: true},
		{PublicKey: record, IsSigner: false, IsWritable: true},
		{PublicKey: domain, IsSigner: false, IsWritable: true},
		{PublicKey: domainOwner, IsSigner: false, IsWritable: true},
		{PublicKey: centralState, IsSigner: false, IsWritable: false},
		{PublicKey: verifier, IsSigner: true, IsWritable: true},
	}
	return solana.NewInstruction(
		programId,
		keys,
		data,
	), nil
}

type WriteRoasInstruction struct {
	Tag   uint8
	RoaId []byte
}

func NewWriteRoaInstruction(roaId []byte) *WriteRoasInstruction {
	return &WriteRoasInstruction{
		Tag:   6,
		RoaId: roaId,
	}
}

func (wri *WriteRoasInstruction) GetInstruction(
	programId,
	systemProgram,
	splNameServiceProgram,
	feePayer,
	record,
	domain,
	domainOwner,
	centralState solana.PublicKey,
) (*solana.GenericInstruction, error) {
	data, err := borsh.Serialize(*wri)
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

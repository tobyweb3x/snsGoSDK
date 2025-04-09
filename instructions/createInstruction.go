package instructions

import (
	"bytes"
	"encoding/binary"
	"slices"

	"github.com/gagliardetto/solana-go"
)

func CreateInstruction(
	nameProgramId,
	systemProgramId,
	nameKey,
	nameOwnerKey,
	payerKey,
	nameClassKey,
	nameParent,
	nameParentOwner solana.PublicKey,
	hashedName []byte, lamports uint64, space uint32) *solana.GenericInstruction {

	var (
		dataBuffer bytes.Buffer
		keys       = make(solana.AccountMetaSlice, 0, 7)
	)
	dataBuffer.WriteByte(0)
	binary.Write(&dataBuffer, binary.LittleEndian, uint32(len(hashedName)))
	dataBuffer.Write(hashedName)
	binary.Write(&dataBuffer, binary.LittleEndian, lamports)
	binary.Write(&dataBuffer, binary.LittleEndian, space)

	keys = append(keys, solana.AccountMetaSlice{
		{PublicKey: systemProgramId, IsSigner: false, IsWritable: false},
		{PublicKey: payerKey, IsSigner: true, IsWritable: true},
		{PublicKey: nameKey, IsSigner: false, IsWritable: true},
		{PublicKey: nameOwnerKey, IsSigner: false, IsWritable: false},
	}...)

	if !nameClassKey.IsZero() {
		keys = append(keys, &solana.AccountMeta{PublicKey: nameClassKey, IsSigner: true, IsWritable: false})
	} else {
		keys = append(keys, &solana.AccountMeta{PublicKey: solana.PublicKey{}, IsSigner: false, IsWritable: false})
	}

	if !nameParent.IsZero() {
		keys = append(keys, &solana.AccountMeta{PublicKey: nameParent, IsSigner: false, IsWritable: false})
	} else {
		keys = append(keys, &solana.AccountMeta{PublicKey: solana.PublicKey{}, IsSigner: false, IsWritable: false})
	}

	if !nameParentOwner.IsZero() {
		keys = append(keys, &solana.AccountMeta{PublicKey: nameParentOwner, IsSigner: true, IsWritable: false})
	}

	return solana.NewInstruction(
		nameProgramId,
		slices.Clip(keys),
		dataBuffer.Bytes(),
	)
}

package instructions

import (
	"bytes"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

func TransferInstruction(
	nameProgramId,
	nameAccountKey,
	newOwnerKey,
	currentOwnerKey,
	nameClassKey,
	nameParent,
	parentOwner solana.PublicKey) *solana.GenericInstruction {
	var dataBuffer bytes.Buffer

	dataBuffer.WriteByte(2)
	binary.Write(&dataBuffer, binary.LittleEndian, newOwnerKey)

	keys := []*solana.AccountMeta{
		{PublicKey: nameAccountKey, IsSigner: false, IsWritable: true},
	}

	if parentOwner.IsZero() {
		keys = append(keys, &solana.AccountMeta{PublicKey: currentOwnerKey, IsSigner: true, IsWritable: false})
	} else {
		keys = append(keys, &solana.AccountMeta{PublicKey: parentOwner, IsSigner: true, IsWritable: false})
	}

	if !nameClassKey.IsZero() {
		keys = append(keys, &solana.AccountMeta{PublicKey: nameClassKey, IsSigner: true, IsWritable: false})
	}

	if !parentOwner.IsZero() && !nameParent.IsZero() {
		if nameClassKey.IsZero() {
			keys = append(keys, &solana.AccountMeta{PublicKey: solana.PublicKeyFromBytes(make([]byte, 32)), IsSigner: false, IsWritable: false})
		}

		keys = append(keys, &solana.AccountMeta{PublicKey: nameParent, IsSigner: false, IsWritable: false})
	}

	return solana.NewInstruction(
		nameProgramId,
		keys,
		dataBuffer.Bytes(),
	)
}

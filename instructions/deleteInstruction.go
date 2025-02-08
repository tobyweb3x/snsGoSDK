package instructions

import (
	"bytes"

	"github.com/gagliardetto/solana-go"
)

func DeleteInstruction(
	nameProgramId, nameAccountKey, refundTargetKey, nameOwnerKey solana.PublicKey,
) *solana.GenericInstruction {

	var dataBuffer bytes.Buffer
	dataBuffer.WriteByte(3)

	keys := []*solana.AccountMeta{
		{PublicKey: nameAccountKey, IsSigner: false, IsWritable: true},
		{PublicKey: nameOwnerKey, IsSigner: true, IsWritable: false},
		{PublicKey: refundTargetKey, IsSigner: false, IsWritable: true},
	}

	return solana.NewInstruction(
		nameProgramId,
		keys,
		dataBuffer.Bytes(),
	)
}

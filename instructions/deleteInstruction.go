package instructions

import (
	"github.com/gagliardetto/solana-go"
)

func DeleteInstruction(
	nameProgramId, nameAccountKey, refundTargetKey, nameOwnerKey solana.PublicKey,
) *solana.GenericInstruction {

	keys := []*solana.AccountMeta{
		{PublicKey: nameAccountKey, IsSigner: false, IsWritable: true},
		{PublicKey: nameOwnerKey, IsSigner: true, IsWritable: false},
		{PublicKey: refundTargetKey, IsSigner: false, IsWritable: true},
	}

	return solana.NewInstruction(
		nameProgramId,
		keys,
		[]byte{3},
	)
}

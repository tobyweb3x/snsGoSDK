package instructions

import (
	"bytes"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

func ReallacInstruction(nameProgramId, systemProgramId, payerkey, nameAccountKey, nameOwnerKey solana.PublicKey, space uint32) *solana.GenericInstruction {
	var dataBuffer bytes.Buffer
	dataBuffer.WriteByte(4)

	binary.Write(&dataBuffer, binary.LittleEndian, space)

	keys := []*solana.AccountMeta{
		{PublicKey: systemProgramId, IsSigner: false, IsWritable: false},
		{PublicKey: payerkey, IsSigner: true, IsWritable: true},
		{PublicKey: nameAccountKey, IsSigner: false, IsWritable: true},
		{PublicKey: nameOwnerKey, IsSigner: true, IsWritable: false},
	}

	return solana.NewInstruction(
		nameProgramId,
		keys,
		dataBuffer.Bytes(),
	)
}

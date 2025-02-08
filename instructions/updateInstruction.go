package instructions

import (
	"bytes"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

func UpdateInstruction(nameProgramId, nameAccountKey, nameUpdateSigner solana.PublicKey, offset uint32, inputData []byte) *solana.GenericInstruction {
	var dataBuffer bytes.Buffer

	dataBuffer.WriteByte(1)

	binary.Write(&dataBuffer, binary.LittleEndian, offset)

	binary.Write(&dataBuffer, binary.LittleEndian, uint32(len(inputData)))

	dataBuffer.Write(inputData)

	keys := []*solana.AccountMeta{
		{PublicKey: nameAccountKey, IsSigner: false, IsWritable: true},
		{PublicKey: nameUpdateSigner, IsSigner: true, IsWritable: false},
	}

	return solana.NewInstruction(
		nameProgramId,
		keys,
		dataBuffer.Bytes(),
	)
}

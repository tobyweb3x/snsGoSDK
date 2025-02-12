package utils

import (
	"bytes"
	"encoding/binary"
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
)

func GetPythFeedAccountKey(shard int, priceFeed []byte) (solana.PublicKey, uint8, error) {
	var buffer *bytes.Buffer
	if err := binary.Write(buffer, binary.LittleEndian, shard); err != nil {
		return solana.PublicKey{}, 0, err
	}

	return solana.FindProgramAddress(
		[][]byte{
			buffer.Bytes(),
			priceFeed,
		},
		spl.DefaultPythPushProgram,
	)
}

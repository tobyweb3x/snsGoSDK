package utils

import (
	"encoding/binary"
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
)

func GetPythFeedAccountKey(shard uint16, priceFeed []byte) (solana.PublicKey, uint8, error) {
	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, shard)
	return solana.FindProgramAddress(
		[][]byte{
			buffer,
			priceFeed,
		},
		spl.DefaultPythPushProgram,
	)
}

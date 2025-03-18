package record

import (
	"bytes"
	"encoding/hex"
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
)

func SerializeSolRecord(
	content, recordKey, signer solana.PublicKey,
	signature []byte,
) ([]byte, error) {

	var dataBuffer bytes.Buffer
	dataBuffer.Write(content.Bytes())
	dataBuffer.Write(recordKey.Bytes())
	// expected := append(content.Bytes(), recordKey.Bytes()...)
	hex, err := hex.DecodeString(dataBuffer.String())
	if err != nil {
		return nil, err
	}
	encodedMessage := []byte(hex)
	valid := CheckSolRecord(encodedMessage, signature, signer)
	if !valid {
		return nil, spl.NewSNSError(spl.InvalidSignature, "The SOL signature is invalid", nil)

	}
	
	return append(content.Bytes(), signature...), nil
}

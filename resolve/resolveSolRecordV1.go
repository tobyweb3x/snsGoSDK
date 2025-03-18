package resolve

import (
	"bytes"
	"encoding/hex"
	"snsGoSDK/record"
	"snsGoSDK/spl"
	"snsGoSDK/types"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func ResolveSolRecordV1(
	conn *rpc.Client,
	owner solana.PublicKey,
	domain string,
) (solana.PublicKey, error) {

	solRecord, err := record.GetSOLRecordRaw(conn, domain)
	if err != nil || solRecord == nil || solRecord.Data == nil {
		return solana.PublicKey{}, spl.NewSNSError(spl.NoRecordData, "the sol record V1 data is empty", err)
	}

	if len(solRecord.Data) < 32 {
		return solana.PublicKey{}, spl.NewSNSError(spl.InvalidRecordData, "the sol record V1 data content is invalid, length < 32", nil)
	}

	var expectedBuffer bytes.Buffer
	expectedBuffer.Write(solRecord.Data[0:32])

	recordKey, err := record.GetRecordKeySync(domain, types.SOL)
	if err != nil {
		return solana.PublicKey{}, err
	}
	expectedBuffer.Write(recordKey.Bytes())

	expected := []byte(hex.EncodeToString(expectedBuffer.Bytes()))
	valid := record.CheckSolRecord(expected, solRecord.Data[32:], owner)
	if !valid {
		return solana.PublicKey{}, spl.NewSNSError(spl.InvalidSignature, "the SOL record V1 signature is invalid", nil)
	}

	return solana.PublicKeyFromBytes(solRecord.Data[0:32]), nil
}

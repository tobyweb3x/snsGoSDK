package record

import (
	"snsGoSDK/spl"
	"snsGoSDK/types"

	"github.com/gagliardetto/solana-go/rpc"
)

func GetRecordDeserialized(
	conn *rpc.Client,
	domain string,
	record types.Record,
) (string, error) {

	pubkey, err := GetRecordKeySync(domain, record)
	if err != nil {
		return "", err
	}
	var nm spl.NameRegistryState
	out, err := nm.Retrieve(conn, pubkey)
	if err != nil || out.Registry.Data == nil {
		return "", spl.NewSNSError(spl.NoRecordData, "the record data is empty", err)
	}

	return DeserializeRecord(
		nm,
		record,
		pubkey,
	)
}

func GetRecordRaw(
	conn *rpc.Client,
	domain string,
	record types.Record,
) (*spl.NameRegistryState, error) {
	pubkey, err := GetRecordKeySync(domain, record)
	if err != nil {
		return nil, err
	}

	var nm spl.NameRegistryState
	out, err := nm.Retrieve(conn, pubkey)
	if err != nil || out.Registry.Data == nil {
		return nil, spl.NewSNSError(spl.NoRecordData, "the record data is empty", err)
	}

	recordSize, ok := types.RecordV1Size[record]
	if ok && len(out.Registry.Data) > int(recordSize) {
		out.Registry.Data = out.Registry.Data[:recordSize]
	} else {
		return nil, spl.NewSNSError(spl.InvalidRecordData, "the record data content is invalid", nil)
	}

	return out.Registry, nil
}

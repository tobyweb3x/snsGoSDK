package recordv2

import (
	"github.com/gagliardetto/solana-go/rpc"

	snsRecord "snsGoSDK/sns-record"
	"snsGoSDK/types"
)

type GetRecordV2Result struct {
	RetrievedRecord    snsRecord.Record
	Record             types.Record
	DeserializeContent string
}

func GetRecordV2(
	conn *rpc.Client,
	domain string,
	record types.Record,
	deserialize bool,
) (GetRecordV2Result, error) {
	pubkey, err := GetRecordV2Key(domain, record)
	if err != nil {
		return GetRecordV2Result{}, err
	}
	retrieveRecord := snsRecord.Record{}
	if err := retrieveRecord.Retrieve(conn, pubkey); err != nil {
		return GetRecordV2Result{}, err
	}
	if deserialize {
		content, err := retrieveRecord.GetContent()
		if err != nil {
			return GetRecordV2Result{}, err
		}
		deserializeContent, err := DeserializeRecordV2Content(content, record)
		if err != nil {
			return GetRecordV2Result{}, err
		}

		return GetRecordV2Result{
			RetrievedRecord:    retrieveRecord,
			DeserializeContent: deserializeContent,
			Record:             record,
		}, nil
	}
	return GetRecordV2Result{
			RetrievedRecord: retrieveRecord,
			Record:          record,
		},
		nil
}

package recordv2

import (
	"github.com/gagliardetto/solana-go/rpc"

	snsRecord "snsGoSDK/sns-record"
	"snsGoSDK/types"

	"github.com/gagliardetto/solana-go"
)

func GetMultipleRecordsV2(
	conn *rpc.Client,
	domain string,
	records []types.Record,
	deserialize bool,
) ([]GetRecordV2Result, error) {
	pubkeys := make([]solana.PublicKey, len(records))
	for i, record := range records {
		pubkeys[i], _ = GetRecordV2Key(domain, record)
	}
	r := snsRecord.Record{}
	retrievedRecords, err := r.RetrieveBatch(conn, pubkeys)
	if err != nil {
		return nil, err
	}
	results := make([]GetRecordV2Result, len(retrievedRecords))

	if deserialize {
		for i, rr := range retrievedRecords {
			if rr.Data == nil {
				results[i] = GetRecordV2Result{}
				continue
			}
			content, err := rr.GetContent()
			if err != nil {
				results[i] = GetRecordV2Result{}
				continue
			}
			deserializeContent, err := DeserializeRecordV2Content(content, records[i])
			if err != nil {
				results[i] = GetRecordV2Result{}
				continue
			}
			results[i] = GetRecordV2Result{
				RetrievedRecord:    rr,
				Record:             records[i],
				DeserializeContent: deserializeContent,
			}
		}
		return results, nil
	}

	for i, rr := range retrievedRecords {
		results[i] = GetRecordV2Result{
			RetrievedRecord: rr,
			Record:          records[i],
		}
	}

	return results, nil

}

package record

import (
	"errors"
	"snsGoSDK/spl"
	"snsGoSDK/types"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetRecordsDeserialized(
	conn *rpc.Client,
	domain string,
	records []types.Record,
) ([]string, error) {

	pubkeys := make([]solana.PublicKey, len(records))
	for i := 0; i < len(records); i++ {
		pubkey, err := GetRecordKeySync(domain, records[i])
		if err != nil {
			pubkeys[i] = solana.PublicKey{}
			continue
		}
		pubkeys[i] = pubkey
	}

	allEmpty := true
	for _, pk := range pubkeys {
		if !pk.IsZero() {
			allEmpty = false
			break
		}
	}

	if allEmpty {
		return nil, spl.NewSNSError(spl.NoRecordData, "the record data is empty", nil)
	}

	nm := spl.NameRegistryState{}
	registries, err := nm.RetrieveBatch(conn, pubkeys)
	if err != nil {
		return nil, err
	}

	if len(registries) != len(records) || len(records) != len(pubkeys) {
		return nil, errors.New("unexpected length of arrays")
	}

	r := make([]string, len(registries))
	for i := 0; i < len(registries); i++ {
		if registries[i] == nil {
			r[i] = ""
			continue
		}

		str, err := DeserializeRecord(*registries[i], records[i], pubkeys[i])
		if err != nil {
			r[i] = ""
			continue
		}
		r[i] = str
	}
	return r, nil
}

func GetRecordsRaw(
	conn *rpc.Client,
	domain string,
	records []types.Record,
) ([]*spl.NameRegistryState, error) {
	pubkeys := make([]solana.PublicKey, len(records))
	for i := 0; i < len(records); i++ {
		pubkey, err := GetRecordKeySync(domain, records[i])
		if err != nil {
			pubkeys[i] = solana.PublicKey{}
			continue
		}
		pubkeys[i] = pubkey
	}

	allEmpty := true
	for _, pk := range pubkeys {
		if !pk.IsZero() {
			allEmpty = false
			break
		}
	}

	if allEmpty {
		return nil, spl.NewSNSError(spl.NoRecordData, "the record data is empty", nil)
	}

	nm := spl.NameRegistryState{}
	registries, err := nm.RetrieveBatch(conn, pubkeys)
	if err != nil {
		return nil, err
	}

	if len(registries) != len(records) || len(records) != len(pubkeys) {
		return nil, errors.New("unexpected length of arrays")
	}

	// recordSize, ok := types.RecordV1Size[record]
	// if ok && len(out.Registry.Data) > int(recordSize) {
	// 	out.Registry.Data = out.Registry.Data[:recordSize]
	// } else {
	// 	return nil, spl.NewSNSError(spl.InvalidRecordData, "the record data content is invalid", nil)
	// }

	return registries, nil
}

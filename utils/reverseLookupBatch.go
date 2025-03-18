package utils

import (
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func ReverseLookUpBatch(conn *rpc.Client, nameAccounts []solana.PublicKey) ([]string, error) {

	reverseLookupAccounts := make([]solana.PublicKey, 0, len(nameAccounts))
	for _, v := range nameAccounts {
		hashedReverseLookup := GetHashedNameSync(v.String())
		reverseLookupAccount, _, err := GetNameAccountKeySync(hashedReverseLookup, spl.ReverseLookupClass, solana.PublicKey{})
		if err != nil {
			return nil, err
		}
		reverseLookupAccounts = append(reverseLookupAccounts, reverseLookupAccount)
	}

	ns := spl.NameRegistryState{}
	names, err := ns.RetrieveBatch(conn, reverseLookupAccounts)
	if err != nil {
		return nil, err
	}

	container := make([]string, 0, len(names))
	for i, name := range names {
		if name == nil || name.Data == nil {
			container = append(container, "")
			continue
		}

		d, err := DeserializeReverse(names[i].Data, false)
		if err != nil {
			container = append(container, "")
			continue
		}

		container = append(container, d)
	}

	return container, nil
}

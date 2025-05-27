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
	nameRegistryStates, err := ns.RetrieveBatch(conn, reverseLookupAccounts)
	if err != nil {
		return nil, err
	}

	container := make([]string, 0, len(nameRegistryStates))
	for i, nameRegistry := range nameRegistryStates {
		str := ""
		if nameRegistry != nil && nameRegistry.Data != nil {
			str, _ = DeserializeReverse(nameRegistryStates[i].Data, false)
		}

		container = append(container, str)
	}

	return container, nil
}

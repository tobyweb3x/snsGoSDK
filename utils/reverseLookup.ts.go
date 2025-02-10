package utils

import (
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// ReverseLookup is used to perform a reverse look up.
func ReverseLookup(conn *rpc.Client, nameAccount, parent solana.PublicKey) (string, error) {

	reverseLookupAccount, err := GetReverseKeyFromDomainkey(nameAccount, parent)
	if err != nil {
		return "", err
	}

	var nm spl.NameRegistryState
	registry, err := nm.Retrieve(conn, reverseLookupAccount)
	if err != nil {
		return "", spl.NewSNSError(spl.NoAccountData, "The registry data is empty", err)
	}

	return deserializeReverse(registry.Registry.Data, parent.IsZero())
}

package spl_name_services

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// ReverseLookup is used to perform a reverse look up.
func ReverseLookup(conn *rpc.Client, nameAccount, parent solana.PublicKey) (string, error) {

	reverseLookupAccount, err := GetReverseKeyFromDomainkey(nameAccount, parent)
	if err != nil {
		return "", err
	}

	var nm NameRegistryState
	registry, err := nm.Retrieve(conn, reverseLookupAccount)
	if err != nil {
		return "", NewSNSError(NoAccountData, "The registry data is empty", err)
	}

	return deserializeReverse(registry.Registry.Data, parent.IsZero())
}

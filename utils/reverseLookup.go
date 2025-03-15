package utils

import (
	"fmt"
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

	fmt.Println(reverseLookupAccount, nameAccount)

	var nm spl.NameRegistryState
	registry, err := nm.Retrieve(conn, reverseLookupAccount)
	if err != nil {
		return "", spl.NewSNSError(spl.NoAccountData, "The registry data is empty", err)
	}

	return DeserializeReverse(registry.Registry.Data, parent.IsZero())
}

// 7sHRghUCXnWkPXTsYxB7Bt35eERkPohMefJ17aBJJnXu HPjEbJoeS77Qq31tWuS8pZCsY2yHAW2PcpAKBfETuwLa

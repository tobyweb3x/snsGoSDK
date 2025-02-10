package deprecated

import (
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetNameOwner(conn *rpc.Client, nameAccountKey solana.PublicKey) (spl.RetrieveResult, error) {
	if !nameAccountKey.IsZero() {
		return spl.RetrieveResult{}, spl.NewSNSError(spl.AccountDoesNotExist, "The name account deos not exist", nil)
	}
	var nm spl.NameRegistryState
	return nm.Retrieve(conn, nameAccountKey)
}

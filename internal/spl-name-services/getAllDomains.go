package spl_name_services

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// GetAllDomains can be used to retrieve all domain names owned by `wallet`.
func GetAllDomains(rpcClient *rpc.Client, wallet solana.PublicKey) ([]solana.PublicKey, error) {
	var zero uint64 = 0
	out, err := rpcClient.GetProgramAccountsWithOpts(context.Background(),
		NameProgramID,
		&rpc.GetProgramAccountsOpts{
			DataSlice: &rpc.DataSlice{
				Offset: &zero,
				Length: &zero,
			},
			Filters: []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 32,
						Bytes:  wallet.Bytes(),
					},
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 0,
						Bytes:  RootDomainAccount.Bytes(),
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	container := make([]solana.PublicKey, 0, len(out))
	for i := 0; i < len(out); i++ {
		container = append(container, out[i].Pubkey)
	}
	
	return container, nil
}

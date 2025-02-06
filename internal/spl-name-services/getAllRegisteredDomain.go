package spl_name_services

import (
	"context"

	"github.com/gagliardetto/solana-go/rpc"
)

func GetAllRegisteredDomain(rpcClient *rpc.Client) (rpc.GetProgramAccountsResult, error) {
	var thirtyTwo uint64 = 32
	out, err := rpcClient.GetProgramAccountsWithOpts(context.Background(),
		NameProgramID,
		&rpc.GetProgramAccountsOpts{
			DataSlice: &rpc.DataSlice{
				Offset: &thirtyTwo,
				Length: &thirtyTwo,
			},
			Filters: []rpc.RPCFilter{
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

	return out, nil
}

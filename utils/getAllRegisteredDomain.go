package utils

import (
	"context"

	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go/rpc"
)

func GetAllRegisteredDomain(conn *rpc.Client) (rpc.GetProgramAccountsResult, error) {
	var thirtyTwo uint64 = 32
	out, err := conn.GetProgramAccountsWithOpts(context.Background(),
		spl.NameProgramID,
		&rpc.GetProgramAccountsOpts{
			DataSlice: &rpc.DataSlice{
				Offset: &thirtyTwo,
				Length: &thirtyTwo,
			},
			Filters: []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 0,
						Bytes:  spl.RootDomainAccount.Bytes(),
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

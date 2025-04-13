package nft

import (
	"context"
	"errors"
	"snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetRecordFromMint(
	conn *rpc.Client,
	mint solana.PublicKey,
) ([]*rpc.Account, error) {
	result, err := conn.GetProgramAccountsWithOpts(
		context.TODO(),
		spl.NameTokenizerID,
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					DataSize: 1 + 1 + 32 + 32 + 32,
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 0,
						Bytes:  []byte{2},
					},
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 1 + 1 + 32 + 32,
						Bytes:  mint.Bytes(),
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("empty values were returned")
	}

	accounts := make([]*rpc.Account, len(result))
	for i, v := range result {
		if v == nil || v.Account == nil {
			accounts[i] = nil
			continue
		}

		accounts[i] = v.Account
	}

	return accounts, nil

}

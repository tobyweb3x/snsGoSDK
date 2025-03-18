package twitter

import (
	"context"

	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetTwitterRegistryKeyData(conn *rpc.Client, verifiedPubKey solana.PublicKey) ([]byte, error) {

	filteredAccounts, err := conn.GetProgramAccountsWithOpts(
		context.TODO(),
		spl.NameProgramID,
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 0,
						Bytes:  spl.TwitterRootParentRegistryKey.Bytes(),
					},
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 32,
						Bytes:  verifiedPubKey.Bytes(),
					},
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 64,
						Bytes:  solana.PublicKeyFromBytes(make([]byte, 32)).Bytes(),
					},
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	if len(filteredAccounts) > 1 {
		return nil, spl.NewSNSError(spl.MultipleRegistries, "More than 1 accounts were found", nil)
	}

	dataBytes := filteredAccounts[0].Account.Data.GetBinary()
	return dataBytes[spl.NameRegistryStateHeaderLen:], nil
}

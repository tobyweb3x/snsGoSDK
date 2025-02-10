package twitter

import (
	"context"
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

func getTwitterHandleAndRegistryKeyViaFilters(conn *rpc.Client, verifiedPubkey solana.PublicKey) (solana.PublicKey, string, error) {

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
						Offset: 64,
						Bytes:  verifiedPubkey.Bytes(),
					},
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 64,
						Bytes:  spl.TwittwrVerificationAuthority.Bytes(),
					},
				},
			},
		},
	)

	if err != nil {

	}

	for _, account := range filteredAccounts {
		if len(account.Account.Data.GetBinary()) > spl.HEADER_LEN+32 {
			data := account.Account.Data.GetBinary()[spl.HEADER_LEN:]
			var rt *ReverseTwitterRegistryState
			if err = borsh.Deserialize(rt, data); err != nil {
				return solana.PublicKey{}, "", err
			}
			return solana.PublicKeyFromBytes(rt.TwitterRegistryKey[:]), rt.TwitterHandle, nil
		}
	}

	return solana.PublicKey{}, "", spl.NewSNSError(spl.AccountDoesNotExist, "The twitter account does not exist", nil)

}

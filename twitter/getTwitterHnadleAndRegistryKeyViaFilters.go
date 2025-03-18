package twitter

import (
	"context"
	"errors"
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

func GetTwitterHandleAndRegistryKeyViaFilters(conn *rpc.Client, verifiedPubkey solana.PublicKey) (solana.PublicKey, string, error) {

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
						Bytes:  verifiedPubkey.Bytes(),
					},
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 64,
						Bytes:  spl.TwitterVerificationAuthority.Bytes(),
					},
				},
			},
		},
	)

	if err != nil {
		return solana.PublicKey{}, "", err
	}

	if filteredAccounts == nil {
		return solana.PublicKey{}, "", errors.New("empty result")
	}

	for _, account := range filteredAccounts {
		if account == nil || account.Account == nil || account.Account.Data == nil {
			continue
		}
		if data := account.Account.Data.GetBinary(); len(data) > spl.NameRegistryStateHeaderLen+32 {
			data = data[spl.NameRegistryStateHeaderLen:]
			rt := &ReverseTwitterRegistryState{}
			if err = borsh.Deserialize(rt, data); err != nil {
				return solana.PublicKey{}, "", err
			}
			return solana.PublicKeyFromBytes(rt.TwitterRegistryKey[:]), rt.TwitterHandle, nil
		}
	}

	return solana.PublicKey{}, "", spl.NewSNSError(spl.AccountDoesNotExist, "The twitter account does not exist", nil)

}

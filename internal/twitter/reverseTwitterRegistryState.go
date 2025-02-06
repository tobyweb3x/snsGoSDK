package twitter

import (
	"context"

	spl "snsGoSDK/internal/spl-name-services"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type ReverseTwitterRegistryState struct {
	TwitterRegistryKey [32]byte
	TwitterHandle      string
}

func (rt *ReverseTwitterRegistryState) Deserialize(data []byte) error {

	err := borsh.Deserialize(rt, data[spl.HEADER_LEN:])
	if err != nil {
		return err
	}

	return nil 
}

func (rt *ReverseTwitterRegistryState) Retrieve(rpcClient *rpc.Client, reverseTwitterAccountKey solana.PublicKey) error {
	reverseTwitterAccount, err := rpcClient.GetAccountInfoWithOpts(
		context.TODO(),
		reverseTwitterAccountKey,
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentProcessed,
		},
	)

	if err != nil {
		return spl.NewSNSError(spl.InvalidReverseTwitter, "The reverse twitter account was not found", err)
	}

	if err = rt.Deserialize(reverseTwitterAccount.Bytes()); err != nil {
		return err
	}

	return nil
}

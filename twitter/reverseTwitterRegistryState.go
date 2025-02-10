package twitter

import (
	"context"

	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type ReverseTwitterRegistryState struct {
	TwitterRegistryKey [32]byte
	TwitterHandle      string
}

func NewReverseTwitterRegistryState(twitterRegistryKey [32]byte, twitterHandle string) *ReverseTwitterRegistryState {
	return &ReverseTwitterRegistryState{
		TwitterRegistryKey: twitterRegistryKey,
		TwitterHandle:      twitterHandle,
	}
}

func (rt *ReverseTwitterRegistryState) Retrieve(conn *rpc.Client, reverseTwitterAccountKey solana.PublicKey) error {
	reverseTwitterAccount, err := conn.GetAccountInfoWithOpts(
		context.TODO(),
		reverseTwitterAccountKey,
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentProcessed,
		},
	)

	if err != nil {
		return spl.NewSNSError(spl.InvalidReverseTwitter, "The reverse twitter account was not found", err)
	}

	if len(reverseTwitterAccount.Bytes()) == 0 {
		return spl.NewSNSError(spl.InvalidReverseTwitter, "Reverse Twitter account byte data is empty", nil)
	}

	if len(reverseTwitterAccount.Bytes()) < spl.HEADER_LEN {
		return spl.NewSNSError(spl.InvalidReverseTwitter, "data is too short to contain valid registry state", nil)
	}

	return borsh.Deserialize(rt, reverseTwitterAccount.Bytes()[spl.HEADER_LEN:])
}

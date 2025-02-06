package twitter

import (
	spl "snsGoSDK/internal/spl-name-services"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetHandleAndRegistryKey(rpcClient *rpc.Client, verifiedPubkey solana.PublicKey) (solana.PublicKey, string, error) {

	hashedVerifiedPubkey := spl.GetHashedNameSync(verifiedPubkey.String())
	reverseRegistryKey, _, err := spl.GetNameAccountKeySync(
		hashedVerifiedPubkey,
		spl.TwittwrVerificationAuthority,
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return solana.PublicKey{}, "", err
	}

	var rt ReverseTwitterRegistryState
	if err = rt.Retrieve(rpcClient, reverseRegistryKey); err != nil {
		return solana.PublicKey{}, "", err
	}

	return solana.PublicKeyFromBytes(rt.TwitterRegistryKey[:]),
		rt.TwitterHandle,
		nil
}

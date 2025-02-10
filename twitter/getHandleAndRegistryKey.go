package twitter

import (
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetHandleAndRegistryKey(conn *rpc.Client, verifiedPubkey solana.PublicKey) (solana.PublicKey, string, error) {

	hashedVerifiedPubkey := utils.GetHashedNameSync(verifiedPubkey.String())
	reverseRegistryKey, _, err := utils.GetNameAccountKeySync(
		hashedVerifiedPubkey,
		spl.TwittwrVerificationAuthority,
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return solana.PublicKey{}, "", err
	}

	var rt ReverseTwitterRegistryState
	if err = rt.Retrieve(conn, reverseRegistryKey); err != nil {
		return solana.PublicKey{}, "", err
	}

	return solana.PublicKeyFromBytes(rt.TwitterRegistryKey[:]),
		rt.TwitterHandle,
		nil
}

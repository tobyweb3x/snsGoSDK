package twitter

import (
	"context"
	"reflect"
	"strings"

	spl "snsGoSDK/internal/spl-name-services"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/near/borsh-go"
)

var (
	/*
		The "".twitter" TLD.
			TwitterRootParentRegistryKey = common.PublicKeyFromString("4YcexoW3r78zz16J2aqmukBLRwGq6rAvWzJpkYAXqebv")
	*/
	TwitterRootParentRegistryKey = common.PublicKeyFromString("4YcexoW3r78zz16J2aqmukBLRwGq6rAvWzJpkYAXqebv")

	/*
		The ".twitter" TLD authority.
			TwittwrVerificationAuthority   = common.PublicKeyFromString("FvPH7PrVrLGKPfqaf3xJodFTjZriqrAXXLTVWEorTFBi")
	*/
	TwittwrVerificationAuthority = common.PublicKeyFromString("FvPH7PrVrLGKPfqaf3xJodFTjZriqrAXXLTVWEorTFBi")
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

func (rt *ReverseTwitterRegistryState) Retrieve(rpcClient *client.Client, reverseTwitterAccountKey common.PublicKey) error {
	var (
		reverseTwitterAccount client.AccountInfo
		err                   error
	)
	if reverseTwitterAccount, err = rpcClient.GetAccountInfoWithConfig(context.Background(), reverseTwitterAccountKey.ToBase58(), client.GetAccountInfoConfig{Commitment: rpc.CommitmentProcessed}); err != nil {
		return err
	}

	if reflect.ValueOf(reverseTwitterAccount).IsZero() {
		return ErrInvalidReverseTwitter
	}

	if err = rt.Deserialize(reverseTwitterAccount.Data); err != nil {
		return err
	}

	return nil
}

func GetTwitterRegistry(rpcClient *client.Client, twitterHandle string) (spl.RetrieveResult, error) {
	var (
		twitterHandleRegistryKey common.PublicKey
		nm                       spl.NameRegistryState
		res                      spl.RetrieveResult
		err                      error
	)

	twitterHandle = strings.TrimPrefix(twitterHandle, "@")

	hashedTwitterHandle := spl.GetHashedNameSync(twitterHandle)
	if twitterHandleRegistryKey, _, err = spl.GetNameAccountKeySync(hashedTwitterHandle, spl.NoPublickKeyArg, TwitterRootParentRegistryKey); err != nil {
		return spl.RetrieveResult{}, err
	}

	if res, err = nm.Retrieve(rpcClient, twitterHandleRegistryKey); err != nil {
		return spl.RetrieveResult{}, err
	}

	return res, nil
}

func GetHandleAndRegistryKey(rpcClient *client.Client, verifiedPubkey common.PublicKey) (common.PublicKey, string, error) {
	var (
		rt                 ReverseTwitterRegistryState
		reverseRegistryKey common.PublicKey
		err                error
	)
	hashedVerifiedPubkey := spl.GetHashedNameSync(verifiedPubkey.ToBase58())
	if reverseRegistryKey, _, err = spl.GetNameAccountKeySync(hashedVerifiedPubkey, TwittwrVerificationAuthority, TwitterRootParentRegistryKey); err != nil {
		return common.PublicKey{}, "", err
	}

	if err = rt.Retrieve(rpcClient, reverseRegistryKey); err != nil {
		return common.PublicKey{}, "", err
	}

	return common.PublicKeyFromBytes(rt.TwitterRegistryKey[:]), rt.TwitterHandle, nil

}

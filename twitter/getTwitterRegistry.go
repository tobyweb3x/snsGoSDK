package twitter

import (
	"strings"

	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetTwitterRegistry(conn *rpc.Client, twitterHandle string) (spl.RetrieveResult, error) {

	twitterHandle = strings.TrimPrefix(twitterHandle, "@")

	hashedTwitterHandle := utils.GetHashedNameSync(twitterHandle)
	twitterHandleRegistryKey, _, err := utils.GetNameAccountKeySync(
		hashedTwitterHandle,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return spl.RetrieveResult{}, err
	}
	var nm spl.NameRegistryState
	res, err := nm.Retrieve(conn, twitterHandleRegistryKey)
	if err != nil {
		return spl.RetrieveResult{}, err
	}

	return res, nil
}

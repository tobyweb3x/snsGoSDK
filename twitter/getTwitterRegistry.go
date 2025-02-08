package twitter

import (
	"strings"

	spl "snsGoSDK/spl-name-services"

	"github.com/gagliardetto/solana-go/rpc"
)

func GetTwitterRegistry(conn *rpc.Client, twitterHandle string) (spl.RetrieveResult, error) {

	twitterHandle = strings.TrimPrefix(twitterHandle, "@")

	hashedTwitterHandle := spl.GetHashedNameSync(twitterHandle)
	twitterHandleRegistryKey, _, err := spl.GetNameAccountKeySync(
		hashedTwitterHandle,
		spl.NoPublickKeyArg,
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

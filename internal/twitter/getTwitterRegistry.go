package twitter

import (
	"strings"

	"github.com/gagliardetto/solana-go/rpc"
	spl "snsGoSDK/internal/spl-name-services"
)

func GetTwitterRegistry(rpcClient *rpc.Client, twitterHandle string) (spl.RetrieveResult, error) {

	twitterHandle = strings.TrimPrefix(twitterHandle, "@")

	hashedTwitterHandle := spl.GetHashedNameSync(twitterHandle)
	 twitterHandleRegistryKey, _, err := spl.GetNameAccountKeySync(
		hashedTwitterHandle, 
		spl.NoPublickKeyArg, 
		spl.TwitterRootParentRegistryKey,
		); 
	 if err != nil {
		return spl.RetrieveResult{}, err
	}
	var nm spl.NameRegistryState
	res, err := nm.Retrieve(rpcClient, twitterHandleRegistryKey); 
	if err != nil {
		return spl.RetrieveResult{}, err
	}

	return res, nil
}

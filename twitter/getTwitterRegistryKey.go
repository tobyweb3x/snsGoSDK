package twitter

import (
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func GetTwitterRegistryKey(twitterHandle string) (solana.PublicKey, error) {
	hashedTwitterHandle := utils.GetHashedNameSync(twitterHandle)
	out, _, err := utils.GetNameAccountKeySync(
		hashedTwitterHandle,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return out, nil
}

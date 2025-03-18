package twitter

import (
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func GetTwitterRegistryKey(twitterHandle string) (solana.PublicKey, error) {
	out, _, err := utils.GetNameAccountKeySync(
		utils.GetHashedNameSync(twitterHandle),
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return out, nil
}

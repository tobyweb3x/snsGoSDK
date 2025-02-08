package twitter

import (
	spl "snsGoSDK/spl-name-services"

	"github.com/gagliardetto/solana-go"
)

func getTwitterRegistryKey(twitterHandle string) (solana.PublicKey, error) {
	hashedTwitterHandle := spl.GetHashedNameSync(twitterHandle)
	out, _, err := spl.GetNameAccountKeySync(
		hashedTwitterHandle,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return out, nil
}

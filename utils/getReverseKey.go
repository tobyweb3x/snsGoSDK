package utils

import (
	spl "snsGoSDK/spl"
	"snsGoSDK/types"

	"github.com/gagliardetto/solana-go"
)

// GetReverseKey can be used to get the key of the reverse account
func GetReverseKey(domain string, isSub bool) (solana.PublicKey, error) {
	out, err := GetDomainKeySync(domain, types.VersionUnspecified)
	if err != nil {
		return solana.PublicKey{}, err
	}
	hashedReverseLookup := GetHashedNameSync(out.PubKey.String())
	nameParent := solana.PublicKey{}
	if isSub {
		nameParent = out.Parent
	}
	reverseLookupAccount, _, err := GetNameAccountKeySync(
		hashedReverseLookup,
		spl.ReverseLookupClass,
		nameParent,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return reverseLookupAccount, nil
}

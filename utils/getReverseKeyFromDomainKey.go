package utils

import (
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
)

func GetReverseKeyFromDomainkey(domainKey, parent solana.PublicKey) (solana.PublicKey, error) {
	hashedReverseLookup := GetHashedNameSync(domainKey.String())
	reverseLookupAccount, _, err := GetNameAccountKeySync(
		hashedReverseLookup,
		spl.ReverseLookupClass,
		parent)

	if err != nil {
		return solana.PublicKey{}, err
	}
	return reverseLookupAccount, nil
}

package spl_name_services

import "github.com/gagliardetto/solana-go"

func GetReverseKeyFromDomainkey(domainKey, parent solana.PublicKey) (solana.PublicKey, error) {
	hashedReverseLookup := GetHashedNameSync(domainKey.String())
	reverseLookupAccount, _, err := GetNameAccountKeySync(
		hashedReverseLookup,
		REVERSE_LOOKUP_CLASS,
		parent)

	if err != nil {
		return NoPublickKeyArg, err
	}
	return reverseLookupAccount, nil
}

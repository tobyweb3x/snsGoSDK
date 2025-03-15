package utils

import (
	"strings"

	spl "snsGoSDK/spl"
	"snsGoSDK/types"

	"github.com/gagliardetto/solana-go"
)

func deriveSync(name string, parent, classKey solana.PublicKey) (deriveResult, error) {
	if parent.IsZero() {
		parent = spl.RootDomainAccount
	}
	hashed := GetHashedNameSync(name)
	pubKey, _, err := GetNameAccountKeySync(hashed, classKey, parent)
	if err != nil {
		return deriveResult{}, err
	}
	return deriveResult{PubKey: pubKey, Hashed: hashed}, nil
}

// GetDomainKeySync is used to compute the public key of a domain or subdomain
//
//	`domain` The domain to compute the public key (e.g 'bonfida.sol', 'dex.bonfida.sol')
//
//	`record` Optional parameter: If the domain being resolved is a record
func GetDomainKeySync(domain string, record types.RecordVersion) (DomainKeyResult, error) {
	domain = strings.TrimSuffix(domain, ".sol")

	var (
		recordClass solana.PublicKey
		parentKey   deriveResult
		subKey      deriveResult
		result      deriveResult
		err         error
	)

	if record == types.V2 {
		recordClass = spl.CentralStateSNSRecords
	}

	if splitted := strings.Split(domain, "."); len(splitted) == 2 {

		prefix := string([]byte{uint8(record)})
		subDomain := splitted[0]  // e.g dex
		rootDomain := splitted[1] // e.g  bonfida
		sub := prefix + subDomain

		if parentKey, err = deriveSync(rootDomain, solana.PublicKey{}, solana.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		if result, err = deriveSync(sub, parentKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: true, Parent: parentKey.PubKey}, nil

	} else if len(splitted) == 3 && record != types.VersionUnspecified {

		rootDomain := splitted[2]
		subDomain := splitted[1]
		subRecordDomain := splitted[0]

		// Parent key
		if parentKey, err = deriveSync(rootDomain, spl.RootDomainAccount, solana.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		// Sub domain
		if subKey, err = deriveSync("\x00"+subDomain, parentKey.PubKey, solana.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		// Sub record
		recordPrefix := "\x01"
		if record == types.V2 {
			recordPrefix = "\x02"
		}
		if result, err = deriveSync(recordPrefix+subRecordDomain, subKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: true, Parent: parentKey.PubKey, IsSubRecord: true}, nil

	} else if len(splitted) >= 3 {
		return DomainKeyResult{}, spl.NewSNSError(spl.InvalidInput, "The domain is malformed", nil)
	}

	if result, err = deriveSync(domain, spl.RootDomainAccount, solana.PublicKey{}); err != nil {
		return DomainKeyResult{}, err
	}

	return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: false}, nil
}

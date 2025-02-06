package spl_name_services

import (
	"bytes"
	"strings"

	"github.com/gagliardetto/solana-go"
)

func deriveSync(name string, parent, classKey solana.PublicKey) (deriveResult, error) {
	if parent == NoPublickKeyArg {
		parent = RootDomainAccount
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
func GetDomainKeySync(domain string, record RecordVersion) (DomainKeyResult, error) {
	domain = strings.TrimSuffix(domain, ".sol")

	var (
		recordClass solana.PublicKey
		parentKey   deriveResult
		subKey      deriveResult
		result      deriveResult
		condc       []uint8
		err         error
	)

	if record == V2 {
		recordClass = CentralStateSNSRecords
	}

	if splitted := strings.Split(domain, "."); len(splitted) == 2 {

		condc = []uint8{0}
		if (record == V2) || (record == V1) {
			condc = []uint8{uint8(record)}
		}

		prefix := bytes.NewBuffer(condc).String()

		subDomain := splitted[0]  // e.g bonfida
		rootDomain := splitted[1] // .sol

		sub := prefix + subDomain

		if parentKey, err = deriveSync(rootDomain, NoPublickKeyArg, NoPublickKeyArg); err != nil {
			return DomainKeyResult{}, err
		}
		if result, err = deriveSync(sub, parentKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: true, Parent: parentKey.PubKey}, nil

	} else if len(splitted) == 3 && record != 0 {

		rootDomain := splitted[2]      // .sol
		subDomain := splitted[1]       // e.g bonfida
		subRecordDomain := splitted[0] // e.g dex

		// Parent key
		if parentKey, err = deriveSync(rootDomain, RootDomainAccount, NoPublickKeyArg); err != nil {
			return DomainKeyResult{}, err
		}
		// Sub domain
		if subKey, err = deriveSync("\x00"+subDomain, parentKey.PubKey, NoPublickKeyArg); err != nil {
			return DomainKeyResult{}, err
		}
		// Sub record
		recordPrefix := "\x01"
		if record == V2 {
			recordPrefix = "\x02"
		}
		if result, err = deriveSync(recordPrefix+subRecordDomain, subKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: true, Parent: parentKey.PubKey, IsSubRecord: true}, nil

	} else if len(splitted) >= 3 {
		return DomainKeyResult{}, NewSNSError(InvalidInput, "The domain is malformed", nil)
	}

	if result, err = deriveSync(domain, RootDomainAccount, NoPublickKeyArg); err != nil {
		return DomainKeyResult{}, err
	}

	return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: false}, nil
}

package main

import (
	"crypto/sha256"
	"strings"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
)

type RecordVersion int8

const (
	RecordVersion1 RecordVersion = 1
	RecordVersion2 RecordVersion = 2
)
const (
	HashPrefix string = "SPL Name Service"
)

var (
	NameProgramID          = common.PublicKeyFromString("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")
	CentralStateSNSRecords = common.PublicKeyFromString("2pMnqHvei2N5oDcVGCRdZx48gqt i199v5CsyTTafsbo")
	RootDomainAccount      = common.PublicKeyFromString("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx")
)

func createConnection(connType string) *client.Client {
	return client.NewClient(connType)
}
func GetHashedNameSync(name string) []byte {
	input := HashPrefix + name
	hashed := sha256.Sum256([]byte(input))
	return hashed[:]

}

func GetNameAccountKeySync(hashed256Name []byte, nameClass, nameParent common.PublicKey) (common.PublicKey, uint8, error) {
	var seeds [][]byte
	seeds = append(seeds, hashed256Name)

	if nameClass != (common.PublicKey{}) {
		seeds = append(seeds, nameClass[:])
	} else {
		seeds = append(seeds, make([]byte, 32))
	}

	if nameParent != (common.PublicKey{}) {
		seeds = append(seeds, nameParent[:])
	} else {
		seeds = append(seeds, make([]byte, 32))
	}

	nameAccountKey, nonce, err := common.FindProgramAddress(seeds, NameProgramID)
	if err != nil {
		return common.PublicKey{}, nonce, err

	}

	return nameAccountKey, nonce, nil
}

type DeriveResult struct {
	PubKey common.PublicKey
	Hashed []byte
}

func DeriveSync(name string, parent, classKey common.PublicKey) (DeriveResult, error) {
	if parent != (common.PublicKey{}) {
		parent = RootDomainAccount
	}
	hashed := GetHashedNameSync(name)
	pubKey, _, err := GetNameAccountKeySync(hashed, classKey, parent)
	if err != nil {
		return DeriveResult{}, err
	}
	return DeriveResult{PubKey: pubKey, Hashed: hashed}, nil
}

type DomainKeyResult struct {
	DeriveResult
	IsSub       bool
	Parent      common.PublicKey
	IsSubRecord bool
}

func GetDomainKeySync(domain string, record RecordVersion) (DomainKeyResult, error) {

	domain = strings.TrimSuffix(domain, ".sol")

	var (
		recordClass common.PublicKey
		parentKey   DeriveResult
		subKey      DeriveResult
		result      DeriveResult
		err         error
	)

	if record == RecordVersion2 {
		recordClass = CentralStateSNSRecords
	}

	if splitted := strings.Split(domain, "."); len(splitted) == 2 {

		prefix := []byte{0}
		if record != 0 {
			prefix[0] = byte(record)
		}

		subDomain := splitted[0]
		rootDomain := splitted[1]

		sub := append(prefix, []byte(subDomain)...)

		if parentKey, err = DeriveSync(rootDomain, RootDomainAccount, common.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		if result, err = DeriveSync(string(sub), parentKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{DeriveResult: result, IsSub: true, Parent: parentKey.PubKey}, nil

	} else if len(splitted) == 3 && record != 0 {

		rootDomain := splitted[2]
		subDomain := splitted[1]
		subRecordDomain := splitted[0]

		// Parent key
		if parentKey, err = DeriveSync(rootDomain, RootDomainAccount, common.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		// Sub domain
		if subKey, err = DeriveSync("\x00"+subDomain, parentKey.PubKey, common.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		// Sub record
		recordPrefix := "\x01"
		if record == RecordVersion2 {
			recordPrefix = "\x02"
		}
		if result, err = DeriveSync(recordPrefix+subRecordDomain, subKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{DeriveResult: result, IsSub: true, Parent: parentKey.PubKey, IsSubRecord: true}, nil

	} else if len(splitted) >= 3 {
		return DomainKeyResult{}, ErrInvalidInput
	}

	if result, err = DeriveSync(domain, RootDomainAccount, common.PublicKey{}); err != nil {
		return DomainKeyResult{}, err
	}
	return DomainKeyResult{DeriveResult: result, IsSub: false}, nil
}

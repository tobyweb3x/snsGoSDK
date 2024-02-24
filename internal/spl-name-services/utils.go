package spl_name_services

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
)

func GetDomainKeySync(domain string, record RecordVersion) (DomainKeyResult, error) {

	domain = strings.TrimSuffix(domain, ".sol")

	var (
		recordClass common.PublicKey
		parentKey   deriveResult
		subKey      deriveResult
		result      deriveResult
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

		if parentKey, err = deriveSync(rootDomain, RootDomainAccount, common.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		if result, err = deriveSync(string(sub), parentKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: true, Parent: parentKey.PubKey}, nil

	} else if len(splitted) == 3 && record != 0 {

		rootDomain := splitted[2]
		subDomain := splitted[1]
		subRecordDomain := splitted[0]

		// Parent key
		if parentKey, err = deriveSync(rootDomain, RootDomainAccount, common.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		// Sub domain
		if subKey, err = deriveSync("\x00"+subDomain, parentKey.PubKey, common.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		// Sub record
		recordPrefix := "\x01"
		if record == RecordVersion2 {
			recordPrefix = "\x02"
		}
		if result, err = deriveSync(recordPrefix+subRecordDomain, subKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: true, Parent: parentKey.PubKey, IsSubRecord: true}, nil

	} else if len(splitted) >= 3 {
		return DomainKeyResult{}, ErrInvalidInput
	}

	if result, err = deriveSync(domain, RootDomainAccount, common.PublicKey{}); err != nil {
		return DomainKeyResult{}, err
	}

	return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: false}, nil
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

func deriveSync(name string, parent, classKey common.PublicKey) (deriveResult, error) {
	if parent != (common.PublicKey{}) {
		parent = RootDomainAccount
	}
	hashed := GetHashedNameSync(name)
	pubKey, _, err := GetNameAccountKeySync(hashed, classKey, parent)
	if err != nil {
		return deriveResult{}, err
	}
	return deriveResult{PubKey: pubKey, Hashed: hashed}, nil
}

func ReverseLookup(rpcClient *client.Client, nameAccount common.PublicKey) (string, error) {
	var (
		reverseLookupAccount common.PublicKey
		registry             RetrieveResult
		err                  error
	)

	hashedReverseLookup := GetHashedNameSync(nameAccount.ToBase58())
	if reverseLookupAccount, _, err = GetNameAccountKeySync(hashedReverseLookup, REVERSE_LOOKUP_CLASS, NoPublickKeyArg); err != nil {
		return "", err
	}

	var nm NameRegistryState
	registry, err = nm.Retrieve(rpcClient, reverseLookupAccount)
	if err != nil {
		if !errors.Is(err, ErrIgnored) {
			return "", err
		}
	}

	if len(registry.Registry.Data) == 0 || registry.Registry.Data == nil {
		return "", ErrNoAccountData
	}

	return deserializeReverse(registry.Registry.Data)
}

func deserializeReverse(data []byte) (string, error) {

	if len(data) == 0 || data == nil {
		return "", ErrNoAccountData
	}

	if len(data) < 4 {
		return "", errors.New("data length is less than expected (4)")
	}

	nameLength := binary.LittleEndian.Uint32(data[:4])
	if int(nameLength) > len(data[4:]) {
		return "", errors.New("unexpected data length")
	}
	return string(data[4 : 4+nameLength]), nil

}

func IsPublicKeyEqual(a, b common.PublicKey) bool {
	return a.ToBase58() == b.ToBase58()
}

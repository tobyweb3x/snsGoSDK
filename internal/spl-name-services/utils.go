package spl_name_services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/rpc"
)

// GetDomainKeySync is used to compute the public key of a domain or subdomain.
func GetDomainKeySync(domain string, record RecordVersion) (DomainKeyResult, error) {

	domain = strings.TrimSuffix(domain, ".sol")

	var (
		recordClass common.PublicKey
		parentKey   deriveResult
		subKey      deriveResult
		result      deriveResult
		condc       []uint8
		err         error
	)

	if record == RecordVersion2 {
		recordClass = CentralStateSNSRecords
	}

	if splitted := strings.Split(domain, "."); len(splitted) == 2 {

		condc = []uint8{0}
		if record != 0 {
			condc = []uint8{uint8(record)}
		}

		prefix := bytes.NewBuffer(condc).String()

		subDomain := splitted[0]
		rootDomain := splitted[1]

		sub := prefix + subDomain

		if parentKey, err = deriveSync(rootDomain, NoPublickKeyArg, NoPublickKeyArg); err != nil {
			return DomainKeyResult{}, err
		}
		if result, err = deriveSync(sub, parentKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: true, Parent: parentKey.PubKey}, nil

	} else if len(splitted) == 3 && record != 0 {

		rootDomain := splitted[2]
		subDomain := splitted[1]
		subRecordDomain := splitted[0]

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

	if result, err = deriveSync(domain, RootDomainAccount, NoPublickKeyArg); err != nil {
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

	if nameClass == NoPublickKeyArg {
		seeds = append(seeds, make([]byte, 32))
	} else {
		seeds = append(seeds, nameClass.Bytes())
	}

	if nameParent == NoPublickKeyArg {
		seeds = append(seeds, make([]byte, 32))
	} else {
		seeds = append(seeds, nameParent.Bytes())
	}

	nameAccountKey, nonce, err := common.FindProgramAddress(seeds, NameProgramID)
	if err != nil {
		return common.PublicKey{}, nonce, err
	}

	return nameAccountKey, nonce, nil
}

func deriveSync(name string, parent, classKey common.PublicKey) (deriveResult, error) {
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

// ReverseLookup is used to perform a reverse look up.
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

func ReverseLookUpBatch(rpcClient *client.Client, nameAccounts []common.PublicKey) ([]string, error) {

	var (
		reverseLookupAccount common.PublicKey
		names                []NameRegistryState
		err                  error
	)

	reverseLookupAccounts := make([]common.PublicKey, 0, len(nameAccounts))

	for i := 0; i < len(nameAccounts); i++ {
		hashedReverseLookup := GetHashedNameSync(nameAccounts[i].ToBase58())
		if reverseLookupAccount, _, err = GetNameAccountKeySync(hashedReverseLookup, REVERSE_LOOKUP_CLASS, NoPublickKeyArg); err != nil {
			return nil, err
		}
		reverseLookupAccounts = append(reverseLookupAccounts, reverseLookupAccount)
	}

	var nm NameRegistryState
	if names, err = nm.RetrieveBatch(rpcClient, reverseLookupAccounts); err != nil {
		return nil, err
	}

	container := make([]string, 0, len(names))
	for i := 0; i < len(names); i++ {
		if len(names[i].Data) == 0 || names[i].Data == nil {
			container = append(container, "")
			continue
		}

		d, err := deserializeReverse(names[i].Data)
		if err != nil {
			return nil, err
		}

		container = append(container, d)
	}
	return container, nil
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

// GetAllDomains can be used to retrieve all domain names owned by `wallet`.
func GetAllDomains(rpcClient *client.Client, wallet common.PublicKey) ([]common.PublicKey, error) {

	var (
		result rpc.JsonRpcResponse[rpc.GetProgramAccounts]
		err    error
	)
	filter := rpc.GetProgramAccountsConfig{
		Encoding: rpc.AccountEncodingBase64,
		Filters: []rpc.GetProgramAccountsConfigFilter{
			{
				MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
					Offset: 32,
					Bytes:  wallet.ToBase58(),
				},
			},
			{
				MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
					Offset: 0,
					Bytes:  RootDomainAccount.ToBase58(),
				},
			},
		},
	}
	if result, err = rpcClient.RpcClient.GetProgramAccountsWithConfig(context.Background(), NameProgramID.ToBase58(), filter); err != nil {
		return nil, err
	}

	value := result.GetResult()
	length := len(value)
	container := make([]common.PublicKey, 0, length)

	for i := 0; i < length; i++ {
		container = append(container, common.PublicKeyFromString(value[i].Pubkey))
	}

	return container, nil
}

type GetDomainKeysWithReversesResult struct {
	PubKey common.PublicKey
	Domain string
}

// GetDomainKeysWithReverses can be used to retrieve all domain names owned by `wallet` in a human readable format.
func GetDomainKeysWithReverses(conn *client.Client, wallet common.PublicKey) ([]GetDomainKeysWithReversesResult, error) {

	var (
		encodedNameArr []common.PublicKey
		names          []string
		err            error
	)
	if encodedNameArr, err = GetAllDomains(conn, wallet); err != nil {
		return nil, err
	}

	if names, err = ReverseLookUpBatch(conn, encodedNameArr); err != nil {
		return nil, err
	}

	if len(encodedNameArr) != len(names) {
		return nil, errors.New("length of encodedNameArr and names are not equal")
	}
	container := make([]GetDomainKeysWithReversesResult, 0, len(encodedNameArr))
	for i, v := range encodedNameArr {
		a := GetDomainKeysWithReversesResult{
			PubKey: v,
			Domain: names[i],
		}
		container = append(container, a)
	}
	return container, nil
}

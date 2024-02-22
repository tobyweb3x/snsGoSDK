package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/near/borsh-go"
)

type RecordVersion int8

const (
	RecordVersion1 RecordVersion = 1
	RecordVersion2 RecordVersion = 2
)
const (
	HashPrefix      string         = "SPL Name Service"
	HEADER_LEN                     = 96
	NoCommitmentArg rpc.Commitment = ""
)

var (
	NameProgramID          = common.PublicKeyFromString("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")
	CentralStateSNSRecords = common.PublicKeyFromString("2pMnqHvei2N5oDcVGCRdZx48gqt i199v5CsyTTafsbo")
	RootDomainAccount      = common.PublicKeyFromString("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx")
	/*
		NoPublicKeyArg is an alias for:
			common.PublicKey{}
		so you compare with the equality operator
	*/
	NoPublickKeyArg   = common.PublicKey{}
	MINT_PREFIX       = []byte("tokenized_name")
	NAME_TOKENIZER_ID = common.PublicKeyFromString("nftD3vbNkNqfj2Sd3HZwbpw4BxxKWr4AjGb9X38JeZk")

	TOKEN_PROGRAM_ID     = common.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	REVERSE_LOOKUP_CLASS = common.PublicKeyFromString("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")
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

type deriveResult struct {
	PubKey common.PublicKey
	Hashed []byte
}

func DeriveSync(name string, parent, classKey common.PublicKey) (deriveResult, error) {
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

type DomainKeyResult struct {
	PubKey      common.PublicKey
	Parent      common.PublicKey
	Hashed      []byte
	IsSub       bool
	IsSubRecord bool
}

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

		if parentKey, err = DeriveSync(rootDomain, RootDomainAccount, common.PublicKey{}); err != nil {
			return DomainKeyResult{}, err
		}
		if result, err = DeriveSync(string(sub), parentKey.PubKey, recordClass); err != nil {
			return DomainKeyResult{}, err
		}

		return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: true, Parent: parentKey.PubKey}, nil

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

		return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: true, Parent: parentKey.PubKey, IsSubRecord: true}, nil

	} else if len(splitted) >= 3 {
		return DomainKeyResult{}, ErrInvalidInput
	}

	if result, err = DeriveSync(domain, RootDomainAccount, common.PublicKey{}); err != nil {
		return DomainKeyResult{}, err
	}
	return DomainKeyResult{PubKey: result.PubKey, Hashed: result.Hashed, IsSub: false}, nil
}

type RetrieveResult struct {
	Registry *NameRegistryState
	NftOwner common.PublicKey
}
type NameRegistryState struct {
	ParentName,
	Owner,
	Class common.PublicKey
	Data []byte `borsh_skip:"true"`
}

func (ns *NameRegistryState) Deserialize(data []byte) error {
	var schema struct {
		ParentName [32]byte
		Owner      [32]byte
		Class      [32]byte
	}
	if err := borsh.Deserialize(&schema, data); err != nil {
		return err
	}

	ns.ParentName = common.PublicKeyFromBytes(schema.ParentName[:])
	ns.Owner = common.PublicKeyFromBytes(schema.Owner[:])
	ns.Class = common.PublicKeyFromBytes(schema.Class[:])

	if len(data) > HEADER_LEN {
		ns.Data = data[HEADER_LEN:]
	}

	return nil
}

func (ns *NameRegistryState) Retrieve(conn *client.Client, nameAccountKey common.PublicKey) (RetrieveResult, error) {
	var (
		nameAccount client.AccountInfo
		nftOwner    common.PublicKey
		err         error
	)

	if nameAccount, err = conn.GetAccountInfo(context.Background(), nameAccountKey.ToBase58()); err != nil {
		return RetrieveResult{}, err
	}

	if reflect.ValueOf(nameAccount).IsZero() {
		return RetrieveResult{}, ErrAccountDoesNotExist
	}

	if err = ns.Deserialize(nameAccount.Data); err != nil {
		return RetrieveResult{}, err
	}

	if nftOwner, err = RetrieveNftOwner(conn, nameAccountKey); err != nil {
		if errors.Is(err, ErrZeroMintSupply) {
			return RetrieveResult{
				Registry: ns,
			}, fmt.Errorf("error occured but RetrieveResult{Registry: ns} is set, err: %w", err)
		}

		return RetrieveResult{}, err
	}

	return RetrieveResult{
		Registry: ns,
		NftOwner: nftOwner,
	}, nil

}

func RetrieveNftOwner(conn *client.Client, nameAccount common.PublicKey) (common.PublicKey, error) {

	var (
		mint     common.PublicKey
		mintInfo token.MintAccount
		result   rpc.JsonRpcResponse[rpc.GetProgramAccounts]
		err      error
	)
	seeds := [][]byte{
		MINT_PREFIX,
		nameAccount.Bytes(),
	}

	if mint, _, err = common.FindProgramAddress(seeds, NAME_TOKENIZER_ID); err != nil {
		return common.PublicKey{}, err
	}

	if mintInfo, err = getMint(conn, mint, NoCommitmentArg, NoPublickKeyArg); err != nil {
		return common.PublicKey{}, err
	}

	if mintInfo.Supply == 0 {
		fmt.Printf("\n****-------------------mint supply is %d-------------------****\n", mintInfo.Supply)
		// return common.PublicKey{}, ErrZeroMintSupply
	}

	filter := rpc.GetProgramAccountsConfig{
		Filters: []rpc.GetProgramAccountsConfigFilter{
			{
				MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
					Offset: 0,
					Bytes:  mint.ToBase58(),
				},
			},
			{
				MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
					Offset: 64,
					Bytes:  "2",
				},
			},
			{
				DataSize: 165,
			},
		},
	}

	if result, err = conn.RpcClient.GetProgramAccountsWithConfig(context.Background(), TOKEN_PROGRAM_ID.ToBase58(), filter); err != nil {
		return common.PublicKey{}, err
	}

	if len(result.GetResult()) != 1 {
		return common.PublicKey{}, fmt.Errorf("unexpected length")
	}

	if data, ok := result.GetResult()[0].Account.Data.([]byte); ok {
		return common.PublicKeyFromBytes(data[32:64]), nil

	}

	return common.PublicKey{}, fmt.Errorf("unexpected data type")
}

func getMint(rpcClient *client.Client, address common.PublicKey, commitment rpc.Commitment, programId common.PublicKey) (token.MintAccount, error) {
	accountInfo, err := rpcClient.GetAccountInfo(context.Background(), address.ToBase58())
	if err != nil {
		return token.MintAccount{}, err
	}

	if reflect.ValueOf(accountInfo).IsZero() {
		fmt.Println("IT WAS HERE")
		return token.MintAccount{}, TokenAccountNotFoundError
	}

	if programId == NoPublickKeyArg {
		programId = TOKEN_PROGRAM_ID
	}

	if !IsPublicKeyEqual(accountInfo.Owner, programId) {
		return token.MintAccount{}, TokenInvalidAccountOwnerError
	}

	if len(accountInfo.Data) < token.MintAccountSize {
		return token.MintAccount{}, TokenInvalidAccountSizeError
	}

	mintAccount, err := token.MintAccountFromData(accountInfo.Data)
	if err != nil {
		return token.MintAccount{}, err
	}

	return mintAccount, nil
}

func IsPublicKeyEqual(a, b common.PublicKey) bool {
	return a.ToBase58() == b.ToBase58()
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
	if registry, err = nm.Retrieve(rpcClient, reverseLookupAccount); err != nil {
		return "", err
	}

	if len(registry.Registry.Data) == 0 || registry.Registry.Data == nil {
		return "", NoAccountData
	}

	return deserializeReversee(registry.Registry.Data)
}

func deserializeReversee(data []byte) (string, error) {

	if len(data) == 0 || data == nil {
		return "", NoAccountData
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

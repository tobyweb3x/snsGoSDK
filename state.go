package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/near/borsh-go"
)

var (
	NoCommitmentArg rpc.Commitment = ""
)

const (
	HEADER_LEN = 96
)

var (
	/*
		NoPublicKeyArg is an alias for:
		    common.PublicKey{}
		so you compare with the equality operator
	*/
	NoPublickKeyArg   = common.PublicKey{}
	MINT_PREFIX       = []byte("tokenized_name")
	NAME_TOKENIZER_ID = common.PublicKeyFromString("nftD3vbNkNqfj2Sd3HZwbpw4BxxKWr4AjGb9X38JeZk")

	TOKEN_PROGRAM_ID = common.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
)

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
		fmt.Printf("-------------------mint supply is %d-------------------\n", mintInfo.Supply)
		// return common.PublicKey{}, ErrZeroMintSupply
	}

	filter := rpc.GetProgramAccountsConfig{
		Filters: []rpc.GetProgramAccountsConfigFilter{
			// {
			// 	MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
			// 		Offset: 0,
			// 		Bytes:  mint.ToBase58(),
			// 	},
			// },
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

	if result, err := conn.RpcClient.GetProgramAccountsWithConfig(context.Background(), TOKEN_PROGRAM_ID.ToBase58(), filter); err != nil {
		return common.PublicKey{}, fmt.Errorf("%v : result.GetError()::%w", err, result.GetError())
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

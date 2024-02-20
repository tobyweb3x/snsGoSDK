package main

import (
	"context"
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
	HEADER_LEN int = 96
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

type NameRegistryState struct {
	ParentName, Owner, Class common.PublicKey
	Data                     *[]byte
}

func (ns *NameRegistryState) Deserialize(data []byte) error {
	if err := borsh.Deserialize(ns, data); err != nil {
		return err
	}

	if len(data) > int(HEADER_LEN) {
		data = data[HEADER_LEN:]
		ns.Data = &data
	}

	return nil
}

func Retrieve(conn *client.Client, nameAccountKey common.PublicKey) (map[string]interface{}, error) {
	var (
		nameAccount       client.AccountInfo
		err               error
		nameRegistryState *NameRegistryState
	)
	if nameAccount, err = conn.GetAccountInfo(context.Background(), nameAccountKey.ToBase58()); err != nil {
		return nil, err
	}
	if reflect.ValueOf(nameAccount).IsZero() {
		return nil, ErrAccountDoesNotExist
	}

	if err = nameRegistryState.Deserialize(nameAccount.Data); err != nil {
		return nil, err
	}

	if len(nameAccount.Data) > HEADER_LEN {
		slice := nameAccount.Data[HEADER_LEN:]
		nameRegistryState.Data = &slice
	}

	_, err = RetrieveNftOwner(conn, nameAccountKey)

}

func RetrieveNftOwner(conn *client.Client, nameAccount common.PublicKey) (common.PublicKey, error) {

	var (
		mint     common.PublicKey
		mintInfo token.MintAccount
		result rpc.JsonRpcResponse[rpc.GetProgramAccounts]
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
		return common.PublicKey{}, ErrZeroMintSupply
	}

	filter := rpc.GetProgramAccountsConfig{
		Filters: []rpc.GetProgramAccountsConfigFilter{
			rpc.GetProgramAccountsConfigFilter{
				MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
					Offset: 0,
					Bytes: mint.ToBase58(),
				},
			},
			rpc.GetProgramAccountsConfigFilter{
				MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
					Offset: 64,
					Bytes: "2",
				},
				DataSize: 165,
			},
		},
		
	}

	if result, err := conn.RpcClient.GetProgramAccountsWithConfig(context.Background(), TOKEN_PROGRAM_ID.ToBase58(), filter); err != nil {
		return common.PublicKey{}, fmt.Errorf("%v : %w", err, result.GetError().Error())
	}

	if len(result.GetResult()) != 1 {
		return common.PublicKey{}, nil
	}

	if data, ok := result.GetResult()[0].Account.Data.([]byte); ok {
		
	}
	
	

	return common.PublicKey{}, nil
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

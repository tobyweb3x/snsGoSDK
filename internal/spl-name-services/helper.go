package spl_name_services

import (
	"context"
	"reflect"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
)

func getMint(rpcClient *client.Client, address common.PublicKey, commitment rpc.Commitment, programId common.PublicKey) (token.MintAccount, error) {
	mintAccountInfo, err := rpcClient.GetAccountInfo(context.Background(), address.ToBase58())
	if err != nil {
		return token.MintAccount{}, err
	}

	if reflect.ValueOf(mintAccountInfo).IsZero() {
		return token.MintAccount{}, ErrTokenAccountNotFound
	}

	if programId == NoPublickKeyArg {
		programId = TOKEN_PROGRAM_ID
	}

	if !IsPublicKeyEqual(mintAccountInfo.Owner, programId) {
		return token.MintAccount{}, ErrTokenInvalidAccountOwner
	}

	if len(mintAccountInfo.Data) < token.MintAccountSize {
		return token.MintAccount{}, ErrTokenInvalidAccountSize
	}

	mintAccount, err := token.MintAccountFromData(mintAccountInfo.Data)
	if err != nil {
		return token.MintAccount{}, err
	}

	return mintAccount, nil
}

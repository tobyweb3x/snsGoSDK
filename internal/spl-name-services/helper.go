package spl_name_services

import (
	"context"
	"reflect"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
)

// getMint returns the mint account for the given address. This function is the equivalent of the
// solana web3js function `getMint` and `unpackMint` in the solana web3js library.
// The function can returned these name errors: errTokenAccountNotFound, errTokenInvalidAccountOwner, errTokenInvalidAccountSize,
// which can be ignored.
/*
	mint := someFuncCall()

	mintInfo, err := getMint(conn, mint, NoCommitmentArg, NoPublickKeyArg)
	if err != nil {
		if errors.Is(err, ErrTokenAccountNotFound) || 
			errors.Is(err, ErrTokenInvalidAccountOwner) || 
			errors.Is(err, ErrTokenInvalidAccountSize) {
			
				return mint, ErrIgnored
		}
		return common.PublicKey{}, err
	}
*/
// This is due to accounting for the error throwing of Javascript which does not neccesarilly halt program execution.
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

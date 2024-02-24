package spl_name_services

import (
	"context"
	"errors"
	"fmt"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/rpc"
) 

// RetrieveNftOwner returns the publicKey associated with the nameAccount parameter. It returns
// named errors: ErrIgnored and ErrZeroMintSupply, which can be ignored.
/*
ns := someFuncCall()
if nftOwner, err = RetrieveNftOwner(conn, nameAccountKey); err != nil {
		if errors.Is(err, ErrZeroMintSupply) || 
			errors.Is(err, ErrIgnored) {
			return RetrieveResult{
				Registry: ns,
			}, nil
		}
		return RetrieveResult{}, err
	}
*/
// This is due to accounting for the error throwing of Javascript which does not neccesarilly halt program execution.
func RetrieveNftOwner(conn *client.Client, nameAccount common.PublicKey) (common.PublicKey, error) {

	var (
		mint common.PublicKey
		// mintInfo token.MintAccount
		result rpc.JsonRpcResponse[rpc.GetProgramAccounts]
		err    error
	)
	seeds := [][]byte{
		MINT_PREFIX,
		nameAccount.Bytes(),
	}

	if mint, _, err = common.FindProgramAddress(seeds, NAME_TOKENIZER_ID); err != nil {
		return common.PublicKey{}, err
	}

	mintInfo, err := getMint(conn, mint, NoCommitmentArg, NoPublickKeyArg)
	if err != nil {
		if errors.Is(err, ErrTokenAccountNotFound) || errors.Is(err, ErrTokenInvalidAccountOwner) || errors.Is(err, ErrTokenInvalidAccountSize) {
			return mint, ErrIgnored
		}
		return common.PublicKey{}, err
	}

	if mintInfo.Supply == 0 {
		return mint, ErrZeroMintSupply
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
		return common.PublicKey{}, fmt.Errorf("%v: %w", err, ErrIgnored)
	}

	if len(result.GetResult()) != 1 {
		return common.PublicKey{}, fmt.Errorf("unexpected length")
	}

	if data, ok := result.GetResult()[0].Account.Data.([]byte); ok {
		return common.PublicKeyFromBytes(data[32:64]), nil
	}

	return common.PublicKey{}, fmt.Errorf("unexpected data type")
}

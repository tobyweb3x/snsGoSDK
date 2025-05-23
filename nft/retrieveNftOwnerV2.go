package nft

import (
	"context"
	"errors"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
)

func RetrieveNftOwnerV2(
	conn *rpc.Client,
	nameAccount solana.PublicKey,
) (solana.PublicKey, error) {

	mint, _, err := GetDomainMint(nameAccount)
	if err != nil {
		return solana.PublicKey{}, err
	}
	largestAccount, err := conn.GetTokenLargestAccounts(
		context.TODO(),
		mint,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		var outErr *jsonrpc.RPCError
		if errors.As(err, &outErr) && outErr.Code == -32602 { // Mint does not exist
			return solana.PublicKey{}, nil
		}
		return solana.PublicKey{}, err
	}
	if len(largestAccount.Value) == 0 {
		return solana.PublicKey{}, errors.New("no accounts found")
	}

	largestAccountInfo, err := conn.GetAccountInfo(
		context.TODO(),
		largestAccount.Value[0].Address,
	)
	if err != nil || largestAccountInfo == nil {
		return solana.PublicKey{}, err
	}

	var decoded token.Account
	if err := bin.NewBorshDecoder(largestAccountInfo.Value.Data.GetBinary()).Decode(&decoded); err != nil {
		return solana.PublicKey{}, err
	}

	if decoded.Amount == 1 {
		return decoded.Owner, nil
	}

	return solana.PublicKey{}, errors.New("no accounts found")
}

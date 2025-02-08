package spl_name_services

import (
	"context"
	"errors"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
)

func retrieveNftOwnerV2(conn *rpc.Client, nameAccount solana.PublicKey) (solana.PublicKey, error) {

	mint, err := getDomainMint(nameAccount)
	if err != nil {
		return solana.PublicKey{}, err
	}

	out, err := conn.GetTokenLargestAccounts(context.TODO(), mint, rpc.CommitmentConfirmed)
	if err != nil {
		var outErr *jsonrpc.RPCError
		if errors.As(err, &outErr) && outErr.Code == -32602 { // Mint does not exist
			return solana.PublicKey{}, nil
		}
		return solana.PublicKey{}, err
	}

	if len(out.Value) == 0 {
		return solana.PublicKey{}, errors.New("no accounts found")
	}

	largestAccountInfo, err := conn.GetAccountInfo(context.Background(), out.Value[0].Address)
	if err != nil {
		return solana.PublicKey{}, err
	}

	var mintInfo token.Mint
	if err := bin.NewBinDecoder(largestAccountInfo.Bytes()).Decode(&mint); err != nil {
		return solana.PublicKey{}, err
	}

	if mintInfo.Supply == 1 {
		return largestAccountInfo.Value.Owner, nil
	}

	return solana.PublicKey{}, nil
}

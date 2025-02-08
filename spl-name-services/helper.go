package spl_name_services

import (
	"context"
	"fmt"
	"reflect"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

func getMint(conn *rpc.Client, address solana.PublicKey, commitment rpc.CommitmentType, programId solana.PublicKey) (token.Mint, error) {
	resp, err := conn.GetAccountInfoWithOpts(context.Background(), address, &rpc.GetAccountInfoOpts{
		Commitment: commitment,
	})

	if err != nil {
		return token.Mint{}, err
	}

	fmt.Println("could panic")
	if reflect.ValueOf(resp).IsZero() {
		return token.Mint{}, ErrTokenAccountNotFound
	}
	fmt.Println("did not panic")

	if programId.IsZero() {
		programId = solana.TokenProgramID
	}

	if !resp.Value.Owner.Equals(programId) {
		return token.Mint{}, ErrTokenInvalidAccountOwner
	}

	var mint token.Mint
	if err = bin.NewBinDecoder(resp.Bytes()).Decode(&mint); err != nil {
		return token.Mint{}, err
	}

	return mint, nil
}

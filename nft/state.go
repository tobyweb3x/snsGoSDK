package nft

import (
	"context"
	"snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type Tag uint8

const (
	Uninitialized Tag = iota
	CentralState
	ActiveRecord
	InactiveRecord
)
const Len = 1 + 1 + 32 + 32 + 32

type NftRecord struct {
	Tag                         Tag
	Nonce                       uint8
	NameAccount, Owner, NftMint solana.PublicKey
}

func NewNftRecord(
	tag Tag, nonce uint8,
	nameAccount, owner, nftMint solana.PublicKey,
) *NftRecord {
	return &NftRecord{
		Tag:         tag,
		Nonce:       nonce,
		NameAccount: nameAccount,
		Owner:       owner,
		NftMint:     nftMint,
	}
}

func (nr *NftRecord) Retrieve(conn *rpc.Client, key solana.PublicKey) error {
	out, err := conn.GetAccountInfo(context.TODO(), key)
	if err != nil || out == nil || out.Value == nil {
		return spl.NewSNSError(spl.NftRecordNotFound, "Nft record not found", err)
	}

	return borsh.Deserialize(nr, out.Value.Data.GetBinary())
}

func (nr NftRecord) FindKey(nameAccount, programId solana.PublicKey) (solana.PublicKey, uint8, error) {
	return solana.FindProgramAddress(
		[][]byte{[]byte("nft_record"), nameAccount.Bytes()},
		programId,
	)
}

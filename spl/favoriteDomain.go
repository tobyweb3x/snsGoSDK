package spl

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type FavoriteDmain struct {
	Tag         uint8
	NameAccount solana.PublicKey
}

func NewFavoriteDmain(tag uint8, nameAccount solana.PublicKey) *FavoriteDmain {
	return &FavoriteDmain{
		Tag:         tag,
		NameAccount: nameAccount,
	}
}

func (fd *FavoriteDmain) Retrieve(conn *rpc.Client, key solana.PublicKey) error {
	accountInfo, err := conn.GetAccountInfo(context.TODO(), key)
	if err != nil {
		return NewSNSError(FavouriteDomainNotFound, "The favorite domain does not exist", err)
	}

	return borsh.Deserialize(fd, accountInfo.Value.Data.GetBinary())
}

// GetKeySync can be used to derive the key of a favorite domain
func (fd FavoriteDmain) GetKeySync(programId, owner solana.PublicKey) (solana.PublicKey, error) {
	out, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("favourite_domain"),
			owner.Bytes(),
		},
		programId,
	)
	return out, err
}

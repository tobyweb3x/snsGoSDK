package spl_name_services

import (
	"context"
	"errors"
	"reflect"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/near/borsh-go"
)

// RetrieveResult is a helper struct for NameRegistryState.Retrieve.
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

func (ns *NameRegistryState) deserialize(data []byte) error {
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

	if err = ns.deserialize(nameAccount.Data); err != nil {
		return RetrieveResult{}, err
	}

	if nftOwner, err = RetrieveNftOwner(conn, nameAccountKey); err != nil {
		if errors.Is(err, ErrZeroMintSupply) || errors.Is(err, ErrIgnored) {
			return RetrieveResult{
				Registry: ns,
			}, nil
		}
		return RetrieveResult{}, err
	}

	return RetrieveResult{
		Registry: ns,
		NftOwner: nftOwner,
	}, nil

}

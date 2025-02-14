package spl

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

// RetrieveResult is a helper struct for NameRegistryState.Retrieve.
type RetrieveResult struct {
	Registry *NameRegistryState
	NftOwner solana.PublicKey
}

type NameRegistryState struct {
	headerLen                uint8 `borsh_skip:"true"`
	ParentName, Owner, Class solana.PublicKey
	Data                     []byte `borsh_skip:"true"`
}

func (ns *NameRegistryState) HEADER_LEN() uint8 {
	ns.headerLen = 98
	return ns.headerLen
}

func (ns *NameRegistryState) Deserialize(data []byte) error {
	var schema struct {
		ParentName [32]byte
		Owner      [32]byte
		Class      [32]byte
	}
	if err := borsh.Deserialize(&schema, data); err != nil {
		return fmt.Errorf("borsch deserialization error: %w", err)
	}

	ns.ParentName = solana.PublicKeyFromBytes(schema.ParentName[:])
	ns.Owner = solana.PublicKeyFromBytes(schema.Owner[:])
	ns.Class = solana.PublicKeyFromBytes(schema.Class[:])

	if len(data) > HEADER_LEN {
		ns.Data = data[HEADER_LEN:]
	}

	return nil
}

func (ns *NameRegistryState) Retrieve(conn *rpc.Client, nameAccountKey solana.PublicKey) (RetrieveResult, error) {

	nameAccount, err := conn.GetAccountInfo(context.Background(), nameAccountKey)
	if err != nil || nameAccount.Value.Owner.IsZero() {
		return RetrieveResult{}, NewSNSError(AccountDoesNotExist, "The name account does not exist", err)
	}

	if err = ns.Deserialize(nameAccount.Bytes()); err != nil {
		return RetrieveResult{}, err
	}

	// if err := borsh.Deserialize(ns, nameAccount.Value.Data.Bytes()); err != nil {
	// 	return RetrieveResult{},
	// }

	nftOwner, err := retrieveNftOwnerV2(conn, nameAccountKey)
	if err != nil {
		return RetrieveResult{}, err
	}

	return RetrieveResult{
		Registry: ns,
		NftOwner: nftOwner,
	}, nil

}

func (ns *NameRegistryState) RetrieveBatch(conn *rpc.Client, nameAccountKeys []solana.PublicKey) ([]*NameRegistryState, error) {
	const batchSize = 100
	nameAccounts := make([]*NameRegistryState, 0, len(nameAccountKeys))

	for i := 0; i < len(nameAccountKeys); i += batchSize {
		end := i + batchSize
		if end > len(nameAccountKeys) {
			end = len(nameAccountKeys)
		}

		batchKeys := nameAccountKeys[i:end]
		out, err := conn.GetMultipleAccounts(context.Background(), batchKeys...)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(out.Value); i++ {
			if err := ns.Deserialize(out.Value[i].Data.GetBinary()); err != nil {
				nameAccounts[i] = nil
				continue
			}
			nameAccounts[i] = ns
		}

	}

	return nameAccounts, nil
}

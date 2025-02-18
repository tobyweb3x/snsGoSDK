package spl

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

var schema struct {
	ParentName [32]byte
	Owner      [32]byte
	Class      [32]byte
}

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
	if err := borsh.Deserialize(&schema, data); err != nil {
		return fmt.Errorf("borsch deserialization error: %w", err)
	}

	ns.ParentName = solana.PublicKeyFromBytes(schema.ParentName[:])
	ns.Owner = solana.PublicKeyFromBytes(schema.Owner[:])
	ns.Class = solana.PublicKeyFromBytes(schema.Class[:])

	if len(data) > NameRegistryStateHeaderLen {
		ns.Data = data[NameRegistryStateHeaderLen:]
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
	var batchSize = 100
	nameAccounts := make([]*NameRegistryState, len(nameAccountKeys))

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
			value := out.Value[i]
			if value == nil || value.Data == nil {
				nameAccounts[i] = nil
				continue
			}

			var n = &NameRegistryState{}
			if err := n.Deserialize(value.Data.GetBinary()); err != nil {
				nameAccounts[i] = nil
				continue
			}

			nameAccounts[i] = n
		}
	}

	return nameAccounts, nil
}

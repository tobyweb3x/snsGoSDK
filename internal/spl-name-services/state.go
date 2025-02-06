package spl_name_services

import (
	"context"
	"reflect"

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
	headerLen uint8
	ParentName,
	Owner,
	Class solana.PublicKey
	Data []byte `borsh_skip:"true"`
}

func (ns *NameRegistryState) HEADER_LEN() uint8 {
	ns.headerLen = 98
	return ns.headerLen
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

	ns.ParentName = solana.PublicKeyFromBytes(schema.ParentName[:])
	ns.Owner = solana.PublicKeyFromBytes(schema.Owner[:])
	ns.Class = solana.PublicKeyFromBytes(schema.Class[:])

	if len(data) > HEADER_LEN {
		ns.Data = data[HEADER_LEN:]
	}

	return nil
}

func (ns *NameRegistryState) Retrieve(rpcClient *rpc.Client, nameAccountKey solana.PublicKey) (RetrieveResult, error) {

	nameAccount, err := rpcClient.GetAccountInfo(context.Background(), nameAccountKey)
	if err != nil {
		return RetrieveResult{}, err
	}

	if reflect.ValueOf(nameAccount).IsZero() || nameAccount.Value.Owner.IsZero() {
		return RetrieveResult{}, NewSNSError(AccountDoesNotExist, "The name account does not exist", nil)
	}

	if err = ns.deserialize(nameAccount.GetBinary()); err != nil {
		return RetrieveResult{}, err
	}

	nftOwner, err := retrieveNftOwnerV2(rpcClient, nameAccountKey)
	if err != nil {
		return RetrieveResult{}, err
	}

	return RetrieveResult{
		Registry: ns,
		NftOwner: nftOwner,
	}, nil

}

func (ns *NameRegistryState) RetrieveBatch(rpcClient *rpc.Client, nameAccountKeys []solana.PublicKey) ([]*NameRegistryState, error) {
	const batchSize = 100
	nameAccounts := make([]*NameRegistryState, 0, len(nameAccountKeys))

	for i := 0; i < len(nameAccountKeys); i += batchSize {
		end := i + batchSize
		if end > len(nameAccountKeys) {
			end = len(nameAccountKeys)
		}

		batchKeys := nameAccountKeys[i:end]
		out, err := rpcClient.GetMultipleAccounts(context.Background(), batchKeys...)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(out.Value); i++ {
			if err := ns.deserialize(out.Value[i].Data.GetBinary()); err != nil {
				nameAccounts[i] = nil
				continue
			}
			nameAccounts[i] = ns
		}

	}

	return nameAccounts, nil
}

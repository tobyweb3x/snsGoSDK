package spl

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"golang.org/x/sync/errgroup"
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
	ParentName, Owner, Class solana.PublicKey
	Data                     []byte `borsh_skip:"true"`
}

func (ns *NameRegistryState) Deserialize(data []byte) error {
	if err := borsh.Deserialize(ns, data); err != nil {
		return fmt.Errorf("borsch deserialization error for NameRegistryState: %w", err)
	}
	// ns.ParentName = solana.PublicKeyFromBytes(schema.ParentName[:])
	// ns.Owner = solana.PublicKeyFromBytes(schema.Owner[:])
	// ns.Class = solana.PublicKeyFromBytes(schema.Class[:])

	if len(data) > NameRegistryStateHeaderLen {
		ns.Data = data[NameRegistryStateHeaderLen:]
	}

	return nil
}

func (ns *NameRegistryState) Retrieve(conn *rpc.Client, nameAccountKey solana.PublicKey) (RetrieveResult, error) {

	nameAccount, err := conn.GetAccountInfo(context.Background(), nameAccountKey)
	if err != nil || nameAccount.Value == nil || nameAccount.Value.Data == nil {
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

func (ns *NameRegistryState) RetrieveBat(conn *rpc.Client, nameAccountKeys []solana.PublicKey) ([]*NameRegistryState, error) {
	var batchSize = 100
	nameAccounts := make([]*NameRegistryState, len(nameAccountKeys))

	for i := 0; i < len(nameAccountKeys); i += batchSize {
		end := min(i+batchSize, len(nameAccountKeys))

		batchKeys := nameAccountKeys[i:end]
		out, err := conn.GetMultipleAccounts(context.Background(), batchKeys...)
		if err != nil {
			return nil, err
		}

		for i := range len(out.Value) {
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

func (ns *NameRegistryState) RetrieveBatch(conn *rpc.Client, nameAccountKeys []solana.PublicKey) ([]*NameRegistryState, error) {
	var (
		batchSize    = 100
		mutex        = sync.Mutex{}
		nameAccounts = make([]*NameRegistryState, 0, len(nameAccountKeys))
		g, ctx       = errgroup.WithContext(context.Background())
	)

	for i := 0; i < len(nameAccountKeys); i += batchSize {
		end := min(i+batchSize, len(nameAccountKeys))
		batchKeys := nameAccountKeys[i:end]

		g.Go(func() error {
			out, err := conn.GetMultipleAccounts(ctx, batchKeys...)
			if err != nil {
				return err
			}

			if out == nil || out.Value == nil {
				return errors.New("empty result")
			}

			tempResults := make([]*NameRegistryState, len(out.Value))
			for j, value := range out.Value {
				if value == nil || value.Data == nil {
					tempResults[j] = nil
					continue
				}

				var n = &NameRegistryState{}
				if err := n.Deserialize(value.Data.GetBinary()); err != nil {
					tempResults[j] = nil
					continue
				}

				tempResults[j] = n
			}

			// Lock only when appending to shared slice
			mutex.Lock()
			nameAccounts = append(nameAccounts, tempResults...)
			mutex.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return slices.Clip(nameAccounts), nil
}

package spl

import (
	"context"
	"errors"
	"fmt"
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
	if err != nil || nameAccount == nil || nameAccount.Value == nil || nameAccount.Value.Data == nil {
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

// RetrieveBatch fetches multiple name registry entries in batches of 100 using concurrent goroutines.
// The original order is preserved, and nil is returned for any index where retrieval or deserialization fails.
func (ns *NameRegistryState) RetrieveBatch(conn *rpc.Client, nameAccountKeys []solana.PublicKey) ([]*NameRegistryState, error) {
	var (
		batchSize = 100
		mutex     = sync.Mutex{}
		container = make(map[int]*NameRegistryState, len(nameAccountKeys))
		g, ctx    = errgroup.WithContext(context.Background())
	)

	for i := 0; i < len(nameAccountKeys); i += batchSize {
		batchStart := i // for Go version below 1.22
		end := min(batchStart+batchSize, len(nameAccountKeys))
		batchKeys := nameAccountKeys[batchStart:end]

		g.Go(func() error {
			out, err := conn.GetMultipleAccounts(ctx, batchKeys...)
			if err != nil {
				return err
			}

			if out == nil || out.Value == nil {
				return errors.New("empty result from GetMultipleAccounts")
			}

			for j, value := range out.Value {
				var n *NameRegistryState
				if value != nil && value.Data != nil {
					n = &NameRegistryState{}
					if err := n.Deserialize(value.Data.GetBinary()); err != nil {
						n = nil
					}
				}

				mutex.Lock()
				container[batchStart+j] = n
				mutex.Unlock()
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	nameAccounts := make([]*NameRegistryState, len(nameAccountKeys))
	for i := range len(container) {
		nameAccounts[i] = container[i]
	}

	return nameAccounts, nil
}

func (ns *NameRegistryState) retrieveBatch(conn *rpc.Client, nameAccountKeys []solana.PublicKey) ([]*NameRegistryState, error) {
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

package utils

import (
	"context"
	"slices"
	"snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/sync/errgroup"
)

func FindSubdomains(conn *rpc.Client, parentKey solana.PublicKey) ([]string, error) {
	var (
		g, ctx = errgroup.WithContext(context.Background())
		reverse,
		subs rpc.GetProgramAccountsResult
	)

	g.Go(func() error {
		r, err := conn.GetProgramAccountsWithOpts(
			ctx,
			spl.NameProgramID,
			&rpc.GetProgramAccountsOpts{
				Filters: []rpc.RPCFilter{
					{
						Memcmp: &rpc.RPCFilterMemcmp{
							Offset: 0,
							Bytes:  parentKey.Bytes(),
						},
					},
					{
						Memcmp: &rpc.RPCFilterMemcmp{
							Offset: 64,
							Bytes:  spl.ReverseLookupClass.Bytes(),
						},
					},
				},
			},
		)
		if err != nil {
			return err
		}

		reverse = r
		return nil
	})

	g.Go(func() error {
		var zeroUint64Pointer uint64
		s, err := conn.GetProgramAccountsWithOpts(
			ctx,
			spl.NameProgramID,
			&rpc.GetProgramAccountsOpts{
				DataSlice: &rpc.DataSlice{
					Offset: &zeroUint64Pointer,
					Length: &zeroUint64Pointer,
				},
				Filters: []rpc.RPCFilter{
					{
						Memcmp: &rpc.RPCFilterMemcmp{
							Offset: 0,
							Bytes:  parentKey.Bytes(),
						},
					},
				},
			},
		)
		if err != nil {
			return err
		}
		subs = s
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	revserseMap := make(map[string]string, len(reverse))
	for _, v := range reverse {
		if v == nil || v.Account == nil ||
			len(v.Account.Data.GetBinary()) < 96 || v.Pubkey.IsZero() {
			continue
		}

		str, _ := DeserializeReverse(v.Account.Data.GetBinary()[96:], true)
		revserseMap[v.Pubkey.String()] = str
	}

	subSlice := make([]string, 0, len(subs))
	for _, v := range subs {
		if v == nil || v.Pubkey.IsZero() {
			continue
		}
		revKey, _ := GetReverseKeyFromDomainkey(v.Pubkey, parentKey)
		if rev, ok := revserseMap[revKey.String()]; ok {
			subSlice = append(subSlice, rev)
		}
	}

	return slices.Clip(subSlice), nil
}

package utils

import (
	"bytes"
	"context"
	"fmt"
	"snsGoSDK/nft"
	spl "snsGoSDK/spl"
	"strings"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/sync/errgroup"
)

type GetPrimaryDomainResult struct {
	Domain  solana.PublicKey
	Reverse string
	Stale   bool
}

func GetPrimaryDomain(conn *rpc.Client, owner solana.PublicKey) (GetPrimaryDomainResult, error) {

	var fd spl.FavoriteDmain
	favKey, _, err := fd.GetKey(spl.NameOffersID, owner)
	if err != nil {

	}
	if err := fd.Retrieve(conn, favKey); err != nil {

	}

	var nm spl.NameRegistryState
	out, err := nm.Retrieve(conn, fd.NameAccount)
	if err != nil {

	}

	var domainOwner, parent solana.PublicKey
	if !out.Registry.ParentName.Equals(spl.RootDomainAccount) {
		parent = out.Registry.ParentName
	}

	reverse, err := ReverseLookup(
		conn,
		fd.NameAccount,
		parent,
	)
	if err != nil {

	}
	if !out.Registry.ParentName.Equals(spl.RootDomainAccount) {
		parentReverse, err := ReverseLookup(
			conn,
			out.Registry.ParentName,
			solana.PublicKey{},
		)
		if err != nil {

		}
		reverse = strings.Join([]string{reverse, parentReverse}, ".")
	}

	if out.NftOwner.IsZero() {
		domainOwner = out.Registry.Owner
	} else {
		domainOwner = out.NftOwner
	}

	return GetPrimaryDomainResult{
		Domain:  fd.NameAccount,
		Reverse: reverse,
		Stale:   owner.Equals(domainOwner),
	}, nil

}

func GetMultipleFavoriteDomain(conn *rpc.Client, wallets []solana.PublicKey) ([]string, error) {

	result := make([]string, 0, len(wallets))

	favKeys := make([]solana.PublicKey, 0, len(wallets))
	for i := 0; i < len(wallets); i++ {
		var fd spl.FavoriteDmain
		favKey, _, err := fd.GetKey(spl.NameOffersID, wallets[i])
		if err != nil {
			favKeys = append(favKeys, solana.PublicKey{})
			continue
		}
		favKeys = append(favKeys, favKey)
	}

	out, err := conn.GetMultipleAccounts(context.TODO(), favKeys...)
	if err != nil {
		return nil, err
	}

	favDomains := make([]solana.PublicKey, 0, len(out.Value))
	for i := 0; i < len(out.Value); i++ {
		var fd spl.FavoriteDmain
		if err := fd.Retrieve(conn, wallets[i]); err != nil {
			favDomains = append(favDomains, solana.PublicKey{})
			continue
		}
		favDomains = append(favDomains, fd.NameAccount)
	}

	domainInfos, err := conn.GetMultipleAccounts(context.TODO(), favDomains...)
	if err != nil {
		return nil, err
	}
	parentRevKeys := make([]solana.PublicKey, 0, len(domainInfos.Value))
	revKeys := make([]solana.PublicKey, 0, len(domainInfos.Value))
	for i := 0; i < len(domainInfos.Value); i++ {
		v := domainInfos.Value[i]
		var parent solana.PublicKey
		if v != nil && len(v.Data.GetBinary()) > 32 {
			parent = solana.PublicKeyFromBytes(v.Data.GetBinary()[:32])
		} else {
			parent = solana.PublicKey{}
		}

		isSub := v != nil && v.Owner.Equals(spl.NameProgramID) && !parent.Equals(spl.RootDomainAccount)

		if input, err := GetReverseKeyFromDomainkey(parent, solana.PublicKey{}); isSub && err != nil {
			parentRevKeys = append(parentRevKeys, input)
		} else {
			parentRevKeys = append(parentRevKeys, solana.PublicKey{})
		}

		if !isSub {
			parent = solana.PublicKey{}
		}

		if input, err := GetReverseKeyFromDomainkey(favDomains[i], parent); err != nil {
			revKeys = append(revKeys, input)
		} else {
			revKeys = append(revKeys, solana.PublicKey{})
		}
	}

	atas := make([]solana.PublicKey, 0, len(favDomains))
	for i := 0; i < len(favDomains); i++ {
		v := favDomains[i]
		mint, _, err := nft.GetDomainMint(v)
		if err != nil {
			atas = append(atas, solana.PublicKey{})
			continue
		}
		ata, err := spl.GetAssociatedTokenAddressSync(
			mint,
			wallets[i],
			true,
		)
		if err != nil {
			atas = append(atas, solana.PublicKey{})
			continue
		}

		atas = append(atas, ata)
	}

	var (
		revs       = make([]*rpc.Account, 0, len(revKeys))
		tokenAccs  = make([]*rpc.Account, 0, len(atas))
		parentRevs = make([]*rpc.Account, 0, len(parentRevKeys))
	)

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		out, err := conn.GetMultipleAccounts(ctx, revKeys...)
		if err != nil {
			return err
		}
		revs = append(revs, out.Value...)
		return nil
	})
	g.Go(func() error {
		out, err := conn.GetMultipleAccounts(ctx, atas...)
		if err != nil {
			return err
		}
		tokenAccs = append(tokenAccs, out.Value...)
		return nil
	})
	g.Go(func() error {
		out, err := conn.GetMultipleAccounts(ctx, parentRevKeys...)
		if err != nil {
			return err
		}
		parentRevs = append(parentRevs, out.Value...)
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for i := 0; i < len(wallets); i++ {
		var (
			domainInfo       = domainInfos.Value[i]
			rev              = revs[i]
			parentRevAccount = parentRevs[i]
			tokenAcc         = tokenAccs[i]
			parentRev        string
		)
		if domainInfo == nil || rev == nil {
			result = append(result, "")
			continue
		}

		if parentRevAccount != nil && parentRevAccount.Owner.Equals(spl.NameProgramID) {
			str, err := DeserializeReverse(parentRevAccount.Owner.Bytes()[96:], false)
			if err != nil {

			}
			parentRev = fmt.Sprintf("%s.des", str)
		}

		if len(domainInfo.Data.GetBinary()) < 64 {
			result = append(result, "")
			continue
		}

		nativeOwner := solana.PublicKeyFromBytes(domainInfo.Data.GetBinary()[32:64])
		if nativeOwner.Equals(wallets[i]) {
			str, err := DeserializeReverse(rev.Owner.Bytes()[96:], false)
			if err != nil {
				result = append(result, "")
				continue
			}
			result = append(result, str)
			continue
		}

		// Either tokenized or stale
		if tokenAcc == nil {
			result = append(result, "")
			continue
		}

		var decoded token.Account
		if err := bin.NewBorshDecoder(tokenAcc.Data.GetBinary()).Decode(&decoded); err != nil {
			result = append(result, "")
			continue
		}

		if decoded.Amount == 1 {
			var data bytes.Buffer
			data.Write(rev.Data.GetBinary()[96:])
			data.WriteString(parentRev)
			str, err := DeserializeReverse(data.Bytes(), false)
			if err != nil {
				result = append(result, "")
				continue
			}
			result = append(result, str)
			continue
		}

		result = append(result, "")
	}

	return result, nil
}

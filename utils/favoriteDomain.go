package utils

import (
	"bytes"
	"context"
	"errors"
	"snsGoSDK/nft"
	spl "snsGoSDK/spl"
	"strings"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"golang.org/x/sync/errgroup"
)

type GetPrimaryDomainResult struct {
	Domain  solana.PublicKey
	Reverse string
	Stale   bool
}

func GetFavoriteDoamin(conn *rpc.Client, owner solana.PublicKey) (GetPrimaryDomainResult, error) {
	return GetPrimaryDomain(conn, owner)
}

func GetPrimaryDomain(conn *rpc.Client, owner solana.PublicKey) (GetPrimaryDomainResult, error) {

	var fd spl.FavoriteDmain
	favKey, err := fd.GetKeySync(spl.NameOffersID, owner)
	if err != nil {
		return GetPrimaryDomainResult{}, err
	}

	if err := fd.Retrieve(conn, favKey); err != nil {
		return GetPrimaryDomainResult{}, err
	}

	var nm spl.NameRegistryState
	out, err := nm.Retrieve(conn, fd.NameAccount)
	if err != nil {
		return GetPrimaryDomainResult{}, err
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
		return GetPrimaryDomainResult{}, err
	}

	if !out.Registry.ParentName.Equals(spl.RootDomainAccount) {
		parentReverse, err := ReverseLookup(
			conn,
			out.Registry.ParentName,
			solana.PublicKey{},
		)
		if err != nil {
			return GetPrimaryDomainResult{}, err
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
		Reverse: string(bytes.TrimPrefix([]byte(reverse), []byte("\x00"))),
		Stale:   !owner.Equals(domainOwner),
	}, nil

}

// GetMultipleFavoriteDomain can be used to retrieve the favorite domains for multiple wallets, up to a maximum of 100.
// If a wallet does not have a favorite domain, the result will be 'undefined' instead of the human readable domain as a string.
// This function is optimized for network efficiency, making only four RPC calls, three of which are executed in parallel using Promise.all, thereby reducing the overall execution time.
func GetMultipleFavoriteDomain(conn *rpc.Client, wallets []solana.PublicKey) ([]string, error) {

	favKeys := make([]solana.PublicKey, len(wallets))
	for i, v := range wallets {
		var fd spl.FavoriteDmain
		favKeys[i], _ = fd.GetKeySync(spl.NameOffersID, v)
	}

	out, err := conn.GetMultipleAccounts(context.TODO(), favKeys...)
	if err != nil {
		return nil, err
	}

	if out.Value == nil {
		return nil, errors.New("empty result from call to GetMultipleAccounts")
	}

	// TODO: remove, maybe unneccessary
	if len(out.Value) != len(favKeys) {
		return nil, errors.New("incomplete data from call to GetMultipleAccounts")
	}

	favDomains := make([]solana.PublicKey, len(out.Value))
	for i, v := range out.Value {
		if v != nil && v.Data != nil {
			var fd spl.FavoriteDmain
			if err := borsh.Deserialize(&fd, v.Data.GetBinary()); err == nil {
				favDomains[i] = fd.NameAccount
			}
		}
	}

	domainInfos, err := conn.GetMultipleAccounts(context.TODO(), favDomains...)
	if err != nil {
		return nil, err
	}

	if domainInfos.Value == nil {
		return nil, errors.New("empty result from call to GetMultipleAccounts")
	}

	// TODO: remove, maybe unneccessary
	if len(domainInfos.Value) != len(favDomains) {
		return nil, errors.New("incomplete data from call to GetMultipleAccounts, 2")
	}

	parentRevKeys := make([]solana.PublicKey, len(domainInfos.Value))
	revKeys := make([]solana.PublicKey, len(domainInfos.Value))
	for i, v := range domainInfos.Value {
		var parent solana.PublicKey
		if v != nil && len(v.Data.GetBinary()) > 32 {
			parent = solana.PublicKeyFromBytes(v.Data.GetBinary()[:32])
		}

		input, _ := GetReverseKeyFromDomainkey(parent, solana.PublicKey{})
		parentRevKeys[i] = input

		isSub := v != nil && v.Owner.Equals(spl.NameProgramID) && !parent.Equals(spl.RootDomainAccount)
		if !isSub {
			parent = solana.PublicKey{}
		}

		input, _ = GetReverseKeyFromDomainkey(favDomains[i], parent)
		revKeys[i] = input
	}

	atas := make([]solana.PublicKey, len(favDomains))
	for i, v := range favDomains {
		mint, _, err := nft.GetDomainMint(v)
		if err != nil {
			continue
		}
		ata, err := spl.GetAssociatedTokenAddressSync(
			mint,
			wallets[i],
			true,
		)
		if err != nil {
			continue
		}

		atas[i] = ata
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

	result := make([]string, len(wallets))
	for i, wallet := range wallets {
		var (
			domainInfo       = domainInfos.Value[i]
			rev              = revs[i]
			parentRevAccount = parentRevs[i]
			tokenAcc         = tokenAccs[i]
			parentRev        = ""
		)

		if domainInfo == nil || rev == nil {
			continue
		}

		if parentRevAccount != nil && parentRevAccount.Data != nil &&
			parentRevAccount.Owner.Equals(spl.NameProgramID) {
			if data := parentRevAccount.Data.GetBinary(); len(data) >= 96 {
				if str, err := DeserializeReverse(data[96:], false); err == nil && str != "" {
					parentRev += "." + str
				}
			}
		}

		if domainInfo.Data != nil {
			if data := domainInfo.Data.GetBinary(); len(data) >= 64 {
				nativeOwner := solana.PublicKeyFromBytes(data[32:64])
				if nativeOwner.Equals(wallet) && rev.Data != nil {
					if data := rev.Data.GetBinary(); len(data) >= 96 {
						if str, err := DeserializeReverse(data[96:], true); err == nil {
							result[i] = str + parentRev
							continue
						}
					}
				}
			}
		}

		// Either tokenized or stale
		if tokenAcc == nil {
			continue
		}

		// Tokenized
		var decoded token.Account
		if err := bin.NewBorshDecoder(tokenAcc.Data.GetBinary()).Decode(&decoded); err == nil {
			if decoded.Amount == 1 && rev.Data != nil {
				if data := rev.Data.GetBinary(); len(data) >= 96 {
					if str, err := DeserializeReverse(data[96:], false); err == nil {
						result[i] = str + parentRev
						continue
					}
				}
			}
		}
	}

	return result, nil
}

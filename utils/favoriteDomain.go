package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
	favKey, err := fd.GetKey(spl.NameOffersID, owner)
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

	result := make([]string, 0, len(wallets))

	favKeys := make([]solana.PublicKey, 0, len(wallets))
	for _, v := range wallets {
		var fd spl.FavoriteDmain
		favKey, err := fd.GetKey(spl.NameOffersID, v)
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

	if out.Value == nil {
		return nil, errors.New("empty result from call to GetMultipleAccounts")
	}

	if len(out.Value) != len(favKeys) {
		return nil, errors.New("incomplete data from call to GetMultipleAccounts")
	}

	favDomains := make([]solana.PublicKey, 0, len(out.Value))
	for _, v := range out.Value {
		var fd spl.FavoriteDmain
		if v == nil || v.Data == nil {
			favDomains = append(favDomains, solana.PublicKey{})
			continue
		}
		if err := borsh.Deserialize(&fd, v.Data.GetBinary()); err != nil {
			favDomains = append(favDomains, solana.PublicKey{})
			continue
		}
		favDomains = append(favDomains, fd.NameAccount)
	}

	domainInfos, err := conn.GetMultipleAccounts(context.TODO(), favDomains...)
	if err != nil {
		return nil, err
	}

	if domainInfos.Value == nil {
		return nil, errors.New("empty result from call to GetMultipleAccounts")
	}

	if len(domainInfos.Value) != len(favDomains) {
		return nil, errors.New("incomplete data from call to GetMultipleAccounts, 2")
	}

	parentRevKeys := make([]solana.PublicKey, 0, len(domainInfos.Value))
	revKeys := make([]solana.PublicKey, 0, len(domainInfos.Value))
	for i, v := range domainInfos.Value {
		var parent solana.PublicKey
		if v != nil && len(v.Data.GetBinary()) > 32 {
			parent = solana.PublicKeyFromBytes(v.Data.GetBinary()[:32])
		} else {
			parent = solana.PublicKey{}
		}

		input, _ := GetReverseKeyFromDomainkey(parent, solana.PublicKey{})
		parentRevKeys = append(parentRevKeys, input)

		isSub := v != nil && v.Owner.Equals(spl.NameProgramID) && !parent.Equals(spl.RootDomainAccount)
		if !isSub {
			parent = solana.PublicKey{}
		}

		input, _ = GetReverseKeyFromDomainkey(favDomains[i], parent)
		revKeys = append(revKeys, input)

	}

	atas := make([]solana.PublicKey, 0, len(favDomains))
	for i, v := range favDomains {
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
			if len(parentRevAccount.Owner.Bytes()) > 96 {
				str, _ := DeserializeReverse(parentRevAccount.Owner.Bytes()[96:], false)
				parentRev = fmt.Sprintf(".%s", str)
			}
		}

		if len(domainInfo.Data.GetBinary()) < 64 {
			result = append(result, "")
			continue
		}

		nativeOwner := solana.PublicKeyFromBytes(domainInfo.Data.GetBinary()[32:64])
		if nativeOwner.Equals(wallets[i]) {
			if len(rev.Owner.Bytes()) > 96 {
				str, err := DeserializeReverse(rev.Owner.Bytes()[96:], false)
				if err != nil {
					result = append(result, "")
					continue
				}
				result = append(result, str)
				continue
			}
			//  TODO: incomplete
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

		// Tokenized
		if decoded.Amount == 1 {
			var data bytes.Buffer
			if len(rev.Data.GetBinary()) > 96 {
				data.Write(rev.Data.GetBinary()[96:])
			}
			data.WriteString(parentRev)
			str, err := DeserializeReverse(data.Bytes(), false)
			if err != nil {
				result = append(result, "")
				continue
			}
			result = append(result, fmt.Sprintf("%s%s", str, parentRev))
			continue
		}

		// Stale

		result = append(result, "")
	}

	return result, nil
}

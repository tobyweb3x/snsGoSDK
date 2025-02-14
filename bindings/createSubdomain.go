package bindings

import (
	"context"
	"fmt"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func CreateSubdomain(
	conn *rpc.Client,
	subdomain string,
	space uint64,
	owner, feePayer solana.PublicKey,
) ([]*solana.GenericInstruction, error) {

	sub := strings.Split(subdomain, ".")[0]
	if len(sub) == 0 {
		return nil, spl.NewSNSError(spl.InvalidDomain, "The subdomain name is malformed", nil)
	}
	out, err := utils.GetDomainKeySync(subdomain, types.VersionUnspecified)
	if err != nil {
		return nil, err
	}

	lamports, err := conn.GetMinimumBalanceForRentExemption(
		context.TODO(),
		space+uint64(spl.NameRegistryStateHeaderLen),
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return nil, err
	}

	var whoTopay solana.PublicKey
	if !feePayer.IsZero() {
		whoTopay = feePayer
	} else {
		whoTopay = owner
	}

	ixOne, err := CreateNameRegistry(
		conn,
		fmt.Sprintf("\\0%s", sub),
		whoTopay,
		owner,
		solana.PublicKey{},
		out.Parent,
		space,
		lamports,
	)
	if err != nil {
		return nil, err
	}

	ixns := make([]*solana.GenericInstruction, 2)
	ixns = append(ixns, ixOne)

	reverseKey, err := utils.GetReverseKey(subdomain, true)
	if err != nil {
		return nil, err
	}

	if info, err := conn.GetAccountInfo(context.TODO(), reverseKey); err == nil && info != nil {
		ixTwo, err := CreateReverseName(
			fmt.Sprintf("\\0%s", sub),
			out.PubKey,
			whoTopay,
			out.Parent,
			owner,
		)
		if err != nil {
			return nil, err
		}
		ixns = append(ixns, ixTwo)
	}

	return ixns, nil
}

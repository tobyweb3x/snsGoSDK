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

	if space == 0 {
		space = 2_000
	}

	sub := strings.Split(subdomain, ".")[0]
	if sub == "" {
		return nil, spl.NewSNSError(spl.InvalidDomain, "The subdomain name is malformed", nil)
	}

	out, err := utils.GetDomainKeySync(subdomain, types.V0)
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

	whoTopay := feePayer
	if feePayer.IsZero() {
		whoTopay = owner
	}

	ixOne, err := CreateNameRegistry(
		conn,
		fmt.Sprintf("\x00%s", sub),
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

	ixns := make([]*solana.GenericInstruction, 0, 2)
	ixns = append(ixns, ixOne)

	reverseKey, err := utils.GetReverseKey(subdomain, true)
	if err != nil {
		return nil, err
	}

	if info, err := conn.GetAccountInfo(context.TODO(), reverseKey); err != nil || info == nil {
		ixTwo, err := CreateReverseName(
			fmt.Sprintf("\x00%s", sub),
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

package bindings

import (
	"context"
	"slices"
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func RegisterDomainNameV2(
	conn *rpc.Client,
	name string, space uint32,
	buyer,
	buyerTokenAccount,
	mint,
	referrerKey solana.PublicKey,
) ([]*solana.GenericInstruction, error) {

	// basic validation
	if strings.Contains(name, ".") || strings.TrimSpace(name) != name {
		return nil,
			spl.NewSNSError(spl.InvalidDomain, "The domain name is malformed", nil)
	}

	ixns := make([]*solana.GenericInstruction, 0, 2)
	refIdx := slices.IndexFunc(spl.REFERRERS, func(i solana.PublicKey) bool {
		return i.Equals(referrerKey)
	})

	var (
		refTokenAccount solana.PublicKey
		err             error
	)

	if refIdx != -1 && referrerKey.IsZero() {
		if refTokenAccount, err = spl.GetAssociatedTokenAddressSync(
			mint,
			referrerKey,
			true,
		); err != nil {
			return nil, err
		}

		out, _ := conn.GetAccountInfo(
			context.TODO(),
			refTokenAccount,
		)

		if out == nil || out.Value == nil || out.Value.Data == nil {
			ixn := CreateAssociatedTokenAccountIdempotentInstruction(
				buyer,
				refTokenAccount,
				referrerKey,
				mint,
				solana.PublicKey{},
				solana.PublicKey{},
			)

			ixns = append(ixns, ixn)
		}
	}

	vault, err := spl.GetAssociatedTokenAddressSync(
		mint,
		spl.VaultOwner,
		true,
	)
	if err != nil {
		return nil, err
	}

	pythFeed, ok := spl.PYTHPullFeeds[mint.String()]
	if !ok {
		return nil, spl.NewSNSError(spl.PythFeedNotFound, "The Pyth account for the provided mint was not found", nil)
	}

	pythFeedAccount, _, err := utils.GetPythFeedAccountKey(0, pythFeed)
	if err != nil {
		return nil, err
	}

	hashed := utils.GetHashedNameSync(name)
	nameAccount, _, err := utils.GetNameAccountKeySync(
		hashed,
		solana.PublicKey{},
		spl.RootDomainAccount)
	if err != nil {
		return nil, err
	}

	hashedReverseLookup := utils.GetHashedNameSync(nameAccount.String())
	reverseLookupAccount, _, err := utils.GetNameAccountKeySync(
		hashedReverseLookup,
		spl.CentralState,
		solana.PublicKey{})
	if err != nil {
		return nil, err
	}

	derivedState, _, err := solana.FindProgramAddress(
		[][]byte{nameAccount.Bytes()},
		spl.ResgistryProgramID,
	)
	if err != nil {
		return nil, err
	}

	var referrerIdxOpt *uint16
	if refIdx == -1 {
		referrerIdxOpt = nil
	} else {
		rf := uint16(refIdx)
		referrerIdxOpt = &rf
	}

	ixnTwo, err := instructions.NewCreateSplitV2Instruction(
		name,
		space,
		referrerIdxOpt,
	).GetInstruction(
		spl.ResgistryProgramID,
		spl.NameProgramID,
		spl.RootDomainAccount,
		nameAccount,
		reverseLookupAccount,
		solana.SystemProgramID,
		spl.CentralState,
		buyer,
		buyer,
		buyer,
		buyerTokenAccount,
		pythFeedAccount,
		vault,
		solana.TokenProgramID,
		solana.SysVarRentPubkey,
		derivedState,
		refTokenAccount,
	)
	if err != nil {
		return nil, err
	}

	ixns = append(ixns, ixnTwo)
	return ixns, nil
}

package twitter

import (
	"context"
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func CreateVerifiedTwitterRegistry(
	conn *rpc.Client,
	twitterHandle string,
	verifiedPubKey, payerKey solana.PublicKey,
	space uint32) ([]*solana.GenericInstruction, error) {

	hasedTwitterHandle := utils.GetHashedNameSync(twitterHandle)
	twitterHandleRegistryKey, _, err := utils.GetNameAccountKeySync(
		hasedTwitterHandle,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return nil, err
	}

	lamport, err := conn.GetMinimumBalanceForRentExemption(context.TODO(), uint64(space)+uint64(spl.NameRegistryStateHeaderLen), rpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}

	ixnsOne := instructions.CreateInstruction(
		spl.NameProgramID,
		solana.SystemProgramID,
		twitterHandleRegistryKey,
		verifiedPubKey,
		payerKey,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey,
		spl.TwitterVerificationAuthority, // Twitter authority acts as owner of the parent for all user-facing registries
		hasedTwitterHandle,
		lamport,
		space,
	)

	ixnsTwo, err := CreateReverseTwitterRegistry(
		conn,
		twitterHandle,
		twitterHandleRegistryKey,
		verifiedPubKey,
		payerKey,
	)
	if err != nil {
		return nil, err
	}

	ixnsSlice := make([]*solana.GenericInstruction, 0, len(ixnsTwo)+1)
	ixnsSlice = append(ixnsSlice, ixnsOne)
	ixnsSlice = append(ixnsSlice, ixnsTwo...)

	return ixnsSlice, nil

}

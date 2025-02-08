package twitter

import (
	"context"
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl-name-services"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func CreateVerifiedTwitterRegistry(
	conn *rpc.Client,
	twitterHandle string,
	verifiedPubKey, payerKey solana.PublicKey,
	space uint32) ([]solana.Instruction, error) {

	hasedTwitterHandle := spl.GetHashedNameSync(twitterHandle)
	twitterHandleRegistryKey, _, err := spl.GetNameAccountKeySync(
		hasedTwitterHandle,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey)
	if err != nil {

	}

	lamport, err := conn.GetMinimumBalanceForRentExemption(context.TODO(), uint64(space)+uint64(spl.HEADER_LEN), rpc.CommitmentConfirmed)
	if err != nil {

	}

	ixnOne := instructions.CreateInstruction(
		spl.NameProgramID,
		solana.SystemProgramID,
		twitterHandleRegistryKey,
		verifiedPubKey,
		payerKey,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey,
		spl.TwittwrVerificationAuthority, // Twitter authority acts as owner of the parent for all user-facing registries
		hasedTwitterHandle,
		lamport,
		space,
	)

	ixnTwo, err := CreateReverseTwitterRegistry(
		conn,
		twitterHandle,
		twitterHandleRegistryKey,
		verifiedPubKey,
		payerKey,
	)
	if err != nil {
		return nil, err
	}

	var instructionsList []solana.Instruction
	instructionsList = append(instructionsList, ixnOne)
	instructionsList = append(instructionsList, ixnTwo...)

	return instructionsList, nil

}

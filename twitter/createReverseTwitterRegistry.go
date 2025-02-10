package twitter

import (
	"context"
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

func CreateReverseTwitterRegistry(
	conn *rpc.Client,
	twitterHandle string,
	twitterRegistryKey, verifiedPubKey, payerKey solana.PublicKey) ([]solana.Instruction, error) {

	hashedVerifiedPubkey := utils.GetHashedNameSync(verifiedPubKey.String())
	reverseRegistryKey, _, err := utils.GetNameAccountKeySync(
		hashedVerifiedPubkey,
		spl.TwittwrVerificationAuthority,
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return nil, err
	}

	rt := NewReverseTwitterRegistryState([32]byte(twitterRegistryKey.Bytes()), twitterHandle)
	reverseTwitterRegistryStateBuff, err := borsh.Serialize(*rt)
	if err != nil {
		return nil, err
	}

	lamport, err := conn.GetMinimumBalanceForRentExemption(
		context.TODO(),
		uint64(len(reverseTwitterRegistryStateBuff)+spl.HEADER_LEN),
		rpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}

	return []solana.Instruction{
		instructions.CreateInstruction(
			spl.NameProgramID,
			solana.SystemProgramID,
			reverseRegistryKey,
			verifiedPubKey,
			payerKey,
			spl.TwittwrVerificationAuthority, // Twitter authority acts as class for all reverse-lookup registries
			spl.TwitterRootParentRegistryKey, // Reverse registries are also children of the root
			spl.TwittwrVerificationAuthority,
			hashedVerifiedPubkey,
			lamport,
			uint32(len(reverseTwitterRegistryStateBuff)),
		),
		instructions.UpdateInstruction(
			spl.NameProgramID,
			reverseRegistryKey,
			spl.TwittwrVerificationAuthority,
			0,
			reverseTwitterRegistryStateBuff,
		),
	}, nil

}

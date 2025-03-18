package twitter

import (
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func ChangeVerifiedPubKey(conn *rpc.Client, twitterHandle string,
	currentVerifiedPubKey, newVerifiedPubKey, payerKey solana.PublicKey) ([]*solana.GenericInstruction, error) {

	hashedTwitterHandle := utils.GetHashedNameSync(twitterHandle)
	twitterHandleRegistryKey, _, err := utils.GetNameAccountKeySync(
		hashedTwitterHandle,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return nil, err
	}

	// Transfer the user-facing registry ownership
	ixnOne := instructions.TransferInstruction(
		spl.NameProgramID,
		twitterHandleRegistryKey,
		newVerifiedPubKey,
		currentVerifiedPubKey,
		solana.PublicKey{},
		solana.PublicKey{},
		solana.PublicKey{},
	)

	ixnTwo, err := CreateReverseTwitterRegistry(
		conn,
		twitterHandle,
		twitterHandleRegistryKey,
		newVerifiedPubKey,
		payerKey,
	)
	if err != nil {
		return nil, err
	}

	ixnsSlice := make([]*solana.GenericInstruction, 0, len(ixnTwo)+1)
	ixnsSlice = append(ixnsSlice, ixnOne)
	ixnsSlice = append(ixnsSlice, ixnTwo...)

	return ixnsSlice, nil
}

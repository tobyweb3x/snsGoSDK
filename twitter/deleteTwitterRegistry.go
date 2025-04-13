package twitter

import (
	"snsGoSDK/instructions"
	"snsGoSDK/spl"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func DeleteTwitterRegistry(
	twitterHandle string,
	verifiedPubKey solana.PublicKey,
) ([]*solana.GenericInstruction, error) {
	hashedTwitterHandle := utils.GetHashedNameSync(twitterHandle)
	twitterHandleRegistryKey, _, err := utils.GetNameAccountKeySync(
		hashedTwitterHandle,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return nil, err
	}
	hashedVerifiedPubkey := utils.GetHashedNameSync(verifiedPubKey.String())
	reverseRegistryKey, _, err := utils.GetNameAccountKeySync(
		hashedVerifiedPubkey,
		spl.TwitterVerificationAuthority,
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return nil, err
	}

	return []*solana.GenericInstruction{
		// Delete the user facing registry
		instructions.DeleteInstruction(
			spl.NameProgramID,
			twitterHandleRegistryKey,
			verifiedPubKey,
			verifiedPubKey,
		),
		// Delete the reverse registry
		instructions.DeleteInstruction(
			spl.NameProgramID,
			reverseRegistryKey,
			verifiedPubKey,
			verifiedPubKey,
		),
	}, nil
}

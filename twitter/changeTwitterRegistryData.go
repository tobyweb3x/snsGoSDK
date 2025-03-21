package twitter

import (
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

// ChangeTwitterRegistryData overwrites the data that is written in the user facing registry
// Signed by the verified pubkey
//
//	`offset` is the offset at which to write the input data into the NameRegistryData
func ChangeTwitterRegistryData(twitterHandle string, verifiedPubkey solana.PublicKey, offset uint32, inputData []byte) (solana.Instruction, error) {
	hashedTwitterHandle := utils.GetHashedNameSync(twitterHandle)
	twitterHandleRegistryKey, _, err := utils.GetNameAccountKeySync(
		hashedTwitterHandle,
		solana.PublicKey{},
		spl.TwitterRootParentRegistryKey,
	)
	if err != nil {
		return nil, err
	}

	return instructions.UpdateInstruction(
		spl.NameProgramID,
		twitterHandleRegistryKey,
		verifiedPubkey,
		offset,
		inputData,
	), nil
}

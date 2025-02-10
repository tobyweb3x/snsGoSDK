package bindings

import (
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func CreateReverseName(
	name string,
	nameAccount,
	feePayer,
	parentName,
	parentNameOwner solana.PublicKey) (*solana.GenericInstruction, error) {

	hashedReverseLookup := utils.GetHashedNameSync(nameAccount.String())
	reverseLookupAccount, _, err := utils.GetNameAccountKeySync(
		hashedReverseLookup,
		spl.CentralState,
		parentName,
	)
	if err != nil {
		return nil, err
	}

	return instructions.NewCreateReverseInstruction(name).GetInstruction(
		spl.ResgistryProgramID,
		spl.NameProgramID,
		spl.RootDomainAccount,
		reverseLookupAccount,
		solana.SystemProgramID,
		spl.CentralState,
		feePayer,
		solana.SysVarRentPubkey,
		parentName,
		parentNameOwner,
	)
}

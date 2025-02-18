package bindings

import (
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func DeleteNameRegistry(
	conn *rpc.Client,
	name string,
	refundTargetKey, nameClass, mameParent solana.PublicKey,
) (*solana.GenericInstruction, error) {

	hashedNamed := utils.GetHashedNameSync(name)
	nameAccountKey, _, err := utils.GetNameAccountKeySync(
		hashedNamed,
		nameClass,
		mameParent,
	)
	if err != nil {
		return nil, err
	}

	var nameOwner solana.PublicKey
	if !nameClass.IsZero() {
		nameOwner = nameClass
	} else {
		var nm spl.NameRegistryState
		out, err := nm.Retrieve(conn, nameAccountKey)
		if err != nil {
			return nil, err
		}
		nameOwner = out.Registry.Owner
	}

	return instructions.DeleteInstruction(
		spl.NameProgramID,
		nameAccountKey,
		refundTargetKey,
		nameOwner,
	), nil

}

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

package bindings

import (
	"snsGoSDK/instructions"
	"snsGoSDK/spl"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func TransferNameOwnership(
	conn *rpc.Client,
	name string,
	newOwner,
	nameClass,
	nameParent,
	parentOwner solana.PublicKey,
) (*solana.GenericInstruction, error) {

	hashedName := utils.GetHashedNameSync(name)
	nameAccountKey, _, err := utils.GetNameAccountKeySync(
		hashedName,
		nameClass,
		nameParent,
	)
	if err != nil {
		return nil, err
	}

	var currentOwner solana.PublicKey
	if !nameClass.IsZero() {
		currentOwner = nameClass
	} else {
		var nm spl.NameRegistryState
		out, err := nm.Retrieve(conn, nameAccountKey)
		if err != nil {
			return nil, err
		}
		currentOwner = out.Registry.Owner
	}

	return instructions.TransferInstruction(
		spl.NameProgramID,
		nameAccountKey,
		newOwner,
		currentOwner,
		nameClass,
		nameParent,
		parentOwner,
	), nil
}

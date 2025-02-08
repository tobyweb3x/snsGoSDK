package bindings

import (
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl-name-services"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func deleteNameRegistry(
	conn *rpc.Client,
	name string,
	refundTargetKey, nameClass, mameParent solana.PublicKey,
) (*solana.GenericInstruction, error) {

	hashedNamed := spl.GetHashedNameSync(name)
	nameAccountKey, _, err := spl.GetNameAccountKeySync(
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

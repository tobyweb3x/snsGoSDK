package bindings

import (
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func UpdateNameRegistryData(
	conn *rpc.Client,
	name string,
	offset uint32,
	inputData []byte,
	nameClass, nameParent solana.PublicKey,
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
	var signer solana.PublicKey
	if nameClass.IsZero() {
		signer = nameClass
	} else {
		var nm spl.NameRegistryState
		out, err := nm.Retrieve(conn, nameAccountKey)
		if err != nil {
			return nil, err
		}
		signer = out.Registry.Owner
	}

	return instructions.UpdateInstruction(
		spl.NameProgramID,
		nameAccountKey,
		signer,
		offset,
		inputData,
	), nil

}

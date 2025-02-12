package bindings

import (
	"snsGoSDK/instructions"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func TransferSubdomain(
	conn *rpc.Client,
	subdomain string,
	newOwner,
	owner solana.PublicKey,
	isParentOwnerSigner bool,
) (*solana.GenericInstruction, error) {

	out, err := utils.GetDomainKeySync(subdomain, types.VersionUnspecified)
	if err != nil {
		return nil, err
	}
	if out.Parent.IsZero() || !out.IsSub {
		return nil, spl.NewSNSError(spl.InvalidSubdomain, "the subdomain is not valid", nil)
	}
	if owner.IsZero() {
		var nm spl.NameRegistryState
		out, err := nm.Retrieve(conn, out.PubKey)
		if err != nil {
			return nil, err
		}
		owner = out.Registry.Owner
	}

	var nameParent, nameParentOwner solana.PublicKey

	if isParentOwnerSigner {
		nameParent = out.Parent
		var nm spl.NameRegistryState
		out, err := nm.Retrieve(conn, out.Parent)
		if err != nil {
			return nil, err
		}
		nameParentOwner = out.Registry.Owner
	}

	return instructions.TransferInstruction(
		spl.NameProgramID,
		out.PubKey,
		newOwner,
		owner,
		solana.PublicKey{},
		nameParent,
		nameParentOwner,
	), nil
}

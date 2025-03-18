package bindings

import (
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func RegisterFavorite(
	conn *rpc.Client,
	nameAccount, owner solana.PublicKey,
) (*solana.GenericInstruction, error) {

	var nm spl.NameRegistryState
	out, err := nm.Retrieve(conn, nameAccount)
	if err != nil {
		return nil, err
	}

	var parent solana.PublicKey
	if !out.Registry.ParentName.Equals(spl.RootDomainAccount) {
		parent = out.Registry.ParentName
	}

	var fd spl.FavoriteDmain
	favkey, err := fd.GetKeySync(spl.NameOffersID, owner)
	if err != nil {
		return nil, err
	}
	return instructions.NewRegisterFavoriteInstruction().GetInstruction(
		spl.NameOffersID,
		nameAccount,
		favkey,
		owner,
		solana.SystemProgramID,
		parent,
	)

}

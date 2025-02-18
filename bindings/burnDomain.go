package bindings

import (
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func BurnDomain(domain string, owner, target solana.PublicKey) (*solana.GenericInstruction, error) {
	pubkey, err := utils.GetDomainKeySync(domain, types.V2)
	if err != nil {
		return nil, err
	}
	state, _, err := solana.FindProgramAddress([][]byte{pubkey.PubKey.Bytes()}, spl.ResgistryProgramID)
	if err != nil {
		return nil, err
	}
	resellingState, _, err := solana.FindProgramAddress(
		[][]byte{pubkey.PubKey.Bytes(), {1, 1}},
		spl.ResgistryProgramID,
	)
	if err != nil {
		return nil, err
	}

	reverse, err := utils.GetReverseKey(domain, false)
	if err != nil {
		return nil, err
	}

	return instructions.NewBurnInstruction().GetInstruction(
		spl.ResgistryProgramID,
		spl.NameProgramID,
		solana.SystemProgramID,
		pubkey.PubKey,
		reverse,
		resellingState,
		state,
		spl.ReverseLookupClass,
		owner,
		target,
	)
}

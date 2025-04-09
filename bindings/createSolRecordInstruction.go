package bindings

import (
	"context"
	"fmt"
	"snsGoSDK/instructions"
	"snsGoSDK/record"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func CreateSolRecordInstruction(
	conn *rpc.Client,
	domain string,
	content, signer, payer solana.PublicKey,
	signature []byte,
) (*solana.GenericInstruction, error) {

	out, err := utils.GetDomainKeySync(
		fmt.Sprintf("%s.%s", types.SOL, domain),
		types.V0,
	)
	if err != nil {
		return nil, err
	}

	serialized, err := record.SerializeSolRecord(
		content,
		out.PubKey,
		signer,
		signature,
	)
	if err != nil {
		return nil, err
	}

	space := len(serialized)

	lamport, err := conn.GetMinimumBalanceForRentExemption(
		context.TODO(),
		uint64(space+spl.NameRegistryStateHeaderLen),
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return nil, err
	}
	return instructions.CreateInstruction(
		spl.NameProgramID,
		solana.SystemProgramID,
		out.PubKey,
		signer,
		payer,
		solana.PublicKey{},
		out.Parent,
		signer,
		out.Hashed,
		lamport,
		uint32(space),
	), nil

}

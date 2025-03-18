package bindings

import (
	"context"
	"snsGoSDK/deprecated"
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func CreateNameRegistry(
	conn *rpc.Client,
	name string,
	payerKey, nameOwner, nameClass, parentName solana.PublicKey,
	space, lamport uint64,
) (*solana.GenericInstruction, error) {
	hashedName := utils.GetHashedNameSync(name)
	nameAccountKey, _, err := utils.GetNameAccountKeySync(
		hashedName,
		nameClass,
		parentName,
	)
	if err != nil {
		return nil, err
	}

	balance := lamport
	if lamport < 0.0 {
		if balance, err = conn.GetMinimumBalanceForRentExemption(
			context.TODO(),
			space,
			rpc.CommitmentConfirmed,
		); err != nil {
			return nil, err
		}
	}

	var nameParentOwner solana.PublicKey
	if !parentName.IsZero() {
		out, err := deprecated.GetNameOwner(conn, parentName)
		if err != nil {
			return nil, err
		}
		nameParentOwner = out.Registry.Owner
	}

	return instructions.CreateInstruction(
		spl.NameProgramID,
		solana.SystemProgramID,
		nameAccountKey,
		nameOwner,
		payerKey,
		nameClass,
		parentName,
		nameParentOwner,
		hashedName,
		balance,
		uint32(space),
	), nil

}

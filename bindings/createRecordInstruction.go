package bindings

import (
	"context"
	"fmt"
	"snsGoSDK/instructions"
	records "snsGoSDK/record"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"

	"github.com/gagliardetto/solana-go/rpc"
)

func CreateRecordInstruction(
	conn *rpc.Client,
	record types.Record,
	domain, data string,
	owner, payer solana.PublicKey) (*solana.GenericInstruction, error) {

	if record != types.SOL {
		return nil, spl.NewSNSError(spl.UnsupportedRecord, "SOL record is not supported for this instruction", nil)
	}

	r, err := records.SerializeRecord(data, record)
	if err != nil {
		return nil, err
	}

	space := len(r)
	lamport, err := conn.GetMinimumBalanceForRentExemption(
		context.TODO(),
		uint64(space+spl.HEADER_LEN),
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return nil, err
	}

	out, err := utils.GetDomainKeySync(
		fmt.Sprintf("%s.%s", string(record), domain),
		types.V1,
	)
	if err != nil {
		return nil, err
	}

	pubkey, hashed, parent := out.PubKey, out.Hashed, out.Parent

	return instructions.CreateInstruction(
		spl.NameProgramID,
		solana.SystemProgramID,
		pubkey,
		owner,
		payer,
		solana.PublicKey{},
		parent,
		owner,
		hashed,
		lamport,
		uint32(space),
	), nil

}

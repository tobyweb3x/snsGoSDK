package bindings

import (
	"context"
	"fmt"
	"snsGoSDK/instructions"
	snsRecord "snsGoSDK/record"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func UpdateRecordInstruction(
	conn *rpc.Client,
	domain, data string,
	record types.Record,
	owner, payer solana.PublicKey,
) ([]*solana.GenericInstruction, error) {

	if record == types.SOL {
		return nil, spl.NewSNSError(spl.UnsupportedRecord, "SOL record is not supported for this instruction", nil)
	}

	out, err := utils.GetDomainKeySync(
		fmt.Sprintf("%s.%s", string(record), domain),
		types.V1,
	)
	if err != nil {
		return nil, err
	}

	info, err := conn.GetAccountInfo(context.TODO(), out.PubKey)
	if err != nil || info.Value.Data == nil {
		return nil, spl.NewSNSError(spl.AccountDoesNotExist, "the record account does not exist", err)
	}

	serialized, err := snsRecord.SerializeRecord(
		data,
		record,
	)
	if err != nil {
		return nil, err
	}

	if l := info.Value.Data.GetBinary()[96:]; len(l) != len(serialized) {
		// Delete + create until we can realloc accounts
		ixOne := instructions.DeleteInstruction(
			spl.NameProgramID,
			out.PubKey,
			payer,
			owner,
		)
		ixTwo, err := CreateRecordInstruction(
			conn,
			record,
			domain,
			data,
			owner,
			payer,
		)
		if err != nil {
			return nil, err
		}
		return []*solana.GenericInstruction{ixOne, ixTwo}, nil
	}

	return []*solana.GenericInstruction{instructions.UpdateInstruction(
		spl.NameProgramID,
		out.PubKey,
		owner,
		0,
		serialized,
	)}, nil
}

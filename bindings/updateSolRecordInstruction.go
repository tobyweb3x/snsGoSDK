package bindings

import (
	"context"
	"fmt"
	"snsGoSDK/instructions"
	snsRecord "snsGoSDK/record"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func UpdateSolRecordInstruction(
	conn *rpc.Client,
	domain string,
	content, signer, payer solana.PublicKey,
	signature []byte,
) ([]*solana.GenericInstruction, error) {

	out, err := utils.GetDomainKeySync(
		fmt.Sprintf("%s.%s", string(types.SOL), domain),
		types.V1,
	)
	if err != nil {
		return nil, err
	}

	info, err := conn.GetAccountInfo(context.TODO(), out.PubKey)
	if err != nil || info.Value.Data == nil {
		return nil, spl.NewSNSError(spl.AccountDoesNotExist, "the record account does not exist", err)
	}

	if len(info.Bytes()) != 96 {
		ixOne := instructions.DeleteInstruction(
			spl.NameProgramID,
			out.PubKey,
			payer,
			signer,
		)
		ixTwo, err := CreateSolRecordInstruction(
			conn,
			domain,
			content,
			signer,
			payer,
			signature,
		)
		if err != nil {
			return nil, err
		}
		return []*solana.GenericInstruction{ixOne, ixTwo}, nil
	}
	serialized, err := snsRecord.SerializeSolRecord(
		content,
		out.PubKey,
		signer,
		signature,
	)
	if err != nil {
		return nil, err
	}

	return []*solana.GenericInstruction{
		instructions.UpdateInstruction(
			spl.NameProgramID,
			out.PubKey,
			signer,
			0,
			serialized,
		),
	}, nil

}

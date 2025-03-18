package bindings

import (
	"fmt"
	recordv2 "snsGoSDK/record_v2"
	snsRecord "snsGoSDK/sns-record"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func UpdateRecordV2Instruction(
	domain, content string,
	record types.Record,
	owner, payer solana.PublicKey,
) (*solana.GenericInstruction, error) {

	out, err := utils.GetDomainKeySync(
		fmt.Sprintf("%s.%s", string(record), domain),
		types.V2,
	)
	if err != nil {
		return nil, err
	}
	if out.IsSub {
		out2, err := utils.GetDomainKeySync(
			domain, types.VersionUnspecified,
		)
		if err != nil {
			return nil, err
		}
		out.Parent = out2.PubKey
	}

	if out.Parent.IsZero() {
		return nil, spl.NewSNSError(spl.InvalidParrent, "parent could not be found", nil)
	}

	data, err := recordv2.SerializeRecordV2Content(content, record)
	if err != nil {
		return nil, err
	}

	return snsRecord.EditRecord(
		payer,
		out.PubKey,
		out.Parent,
		owner,
		spl.NameProgramID,
		snsRecord.SNSRecordsID,
		fmt.Sprintf("\x02%s", content),
		data,
	)
}

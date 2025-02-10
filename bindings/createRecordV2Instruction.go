package bindings

import (
	"fmt"
	recordv2 "snsGoSDK/record_v2"
	snsRecord "snsGoSDK/sns-record"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func CreateRecordV2Instruction(
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
		out2, err := utils.GetDomainKeySync(domain, types.V2)
		if err != nil {
			return nil, spl.NewSNSError(spl.InvalidParrent, "parent could not be found", nil)
		}
		out.Parent = out2.PubKey
	}

	data, err := recordv2.SerializeRecordv2Contents(content, record)
	if err != nil {
		return nil, err
	}

	return snsRecord.AllocateAndPostRecord(
		payer,
		out.PubKey,
		out.Parent,
		owner,
		spl.NameProgramID,
		snsRecord.SNSRecordsID,
		fmt.Sprintf("x02%s", domain),
		data,
	)
}

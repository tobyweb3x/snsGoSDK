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
		parent, err := utils.GetDomainKeySync(domain, types.V0)
		if err != nil {
			return nil, spl.NewSNSError(spl.InvalidParrent, "parent could not be found", err)
		}
		out.Parent = parent.PubKey
	}

	if out.Parent.IsZero() {
		return nil, spl.NewSNSError(spl.InvalidParrent, "parent could not be found", nil)
	}

	data, err := recordv2.SerializeRecordV2Content(content, record)
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
		fmt.Sprintf("\x02%s", string(record)),
		data,
	)
}

package bindings

import (
	"fmt"
	snsRecords "snsGoSDK/sns-record"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func DeleteRecordV2(
	domain string,
	record types.Record,
	owner, payer solana.PublicKey) (*solana.GenericInstruction, error) {

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

	return snsRecords.DeleteRecord(
		payer,
		out.Parent,
		owner,
		out.PubKey,
		spl.NameProgramID,
		snsRecords.SNSRecordsID,
	)
}

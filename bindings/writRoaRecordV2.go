package bindings

import (
	"fmt"
	snsRecord "snsGoSDK/sns-record"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func WritRoaRecordV2(
	domain string,
	record types.Record,
	owner, payer, roaId solana.PublicKey,
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
			domain, types.V0,
		)
		if err != nil {
			return nil, err
		}
		out.Parent = out2.PubKey
	}

	if out.Parent.IsZero() {
		return nil, spl.NewSNSError(spl.InvalidParrent, "parent could not be found", nil)
	}

	return snsRecord.WriteRoa(
		payer,
		spl.NameProgramID,
		out.PubKey,
		out.Parent,
		owner,
		roaId,
		snsRecord.SNSRecordsID,
	)
}

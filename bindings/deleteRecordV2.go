package bindings

import (
	"fmt"
	"reflect"
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

	var parent utils.DomainKeyResult
	if out.IsSub {
		if parent, err = utils.GetDomainKeySync(domain, types.VersionUnspecified); err != nil || reflect.ValueOf(parent).IsZero() {
			return nil, spl.NewSNSError(spl.InvalidParrent, "parent could not be found", err)
		}
	}

	return snsRecords.DeleteRecord(
		payer,
		parent.Parent,
		owner,
		out.PubKey,
		spl.NameProgramID,
		snsRecords.SNSRecordsID,
	)

}

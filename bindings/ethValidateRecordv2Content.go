package bindings

import (
	"fmt"
	"reflect"
	snsRecord "snsGoSDK/sns-record"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func EthValidateRecordv2Content(
	domain string,
	record types.Record,
	owner, payer solana.PublicKey,
	signature, expectedPubkey []byte,
) (*solana.GenericInstruction, error) {

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
			return nil, err
		}
	}

	return snsRecord.ValidateEthSignature(
		payer,
		out.Parent,
		parent.Parent,
		owner,
		spl.NameProgramID,
		snsRecord.SNSRecordsID,
		snsRecord.Ethereum,
		signature,
		expectedPubkey,
	)
}

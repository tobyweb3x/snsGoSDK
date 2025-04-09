package bindings

import (
	"fmt"
	snsRecord "snsGoSDK/sns-record"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func ValidateRecordV2Content(
	domain string,
	staleness bool,
	record types.Record,
	owner, payer, verifier solana.PublicKey,
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

	return snsRecord.ValidateSolanaSignature(
		payer,
		out.PubKey,
		out.Parent,
		owner,
		verifier,
		spl.NameProgramID,
		snsRecord.SNSRecordsID,
		staleness,
	)
}

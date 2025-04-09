package recordv2

import (
	"fmt"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func GetRecordV2Key(domain string, record types.Record) (solana.PublicKey, error) {
	out, err := utils.GetDomainKeySync(domain, types.V0)
	if err != nil {
		return solana.PublicKey{}, err
	}

	hashed := utils.GetHashedNameSync(
		fmt.Sprintf("\x02%s", string(record)),
	)
	out2, _, err := utils.GetNameAccountKeySync(
		hashed,
		spl.CentralStateSNSRecords,
		out.PubKey,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return out2, nil

}

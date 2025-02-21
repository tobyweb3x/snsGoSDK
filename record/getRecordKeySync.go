package record

import (
	"fmt"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
)

func GetRecordKeySync(domain string, record types.Record) (solana.PublicKey, error) {
	out, err := utils.GetDomainKeySync(
		fmt.Sprintf("%s.%s", string(record), domain),
		types.V1,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return out.PubKey, nil

}

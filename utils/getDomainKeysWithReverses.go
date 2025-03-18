package utils

import (
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type GetDomainKeysWithReversesResult struct {
	PubKey solana.PublicKey
	Domain string
}

// GetDomainKeysWithReverses can be used to retrieve all domain names owned by `wallet` in a human readable format.
func GetDomainKeysWithReverses(conn *rpc.Client, wallet solana.PublicKey) ([]GetDomainKeysWithReversesResult, error) {

	encodedNameArr, err := GetAllDomains(conn, wallet)
	if err != nil {
		return nil, err
	}

	names, err := ReverseLookUpBatch(conn, encodedNameArr)
	if err != nil {
		return nil, err
	}

	if len(encodedNameArr) != len(names) {
		return nil, errors.New("length of encodedNameArr and names are not equal")
	}

	container := make([]GetDomainKeysWithReversesResult, len(encodedNameArr))
	for i, v := range encodedNameArr {
		a := GetDomainKeysWithReversesResult{
			PubKey: v,
			Domain: names[i],
		}
		container[i] =  a
	}
	return container, nil
}

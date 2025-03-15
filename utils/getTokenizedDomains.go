package utils

import (
	"snsGoSDK/nft"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type GetTokenizedDomainsResult struct {
	Key, Mint solana.PublicKey
	Reverse   string
}

func GetTokenizedDomains(
	conn *rpc.Client,
	owner solana.PublicKey,
) ([]GetTokenizedDomainsResult, error) {

	nftRecords, err := nft.RetrieveRecords(conn, owner)
	if err != nil {
		return nil, err
	}

	nameAccounts := make([]solana.PublicKey, len(nftRecords))
	for i, v := range nftRecords {
		nameAccounts[i] = v.NameAccount
	}

	names, err := ReverseLookUpBatch(conn, nameAccounts)
	if err != nil {
		return nil, err
	}

	result := make([]GetTokenizedDomainsResult, len(names))
	for i, v := range names {
		result[i] = GetTokenizedDomainsResult{
			Key:     nftRecords[i].NameAccount,
			Mint:    nftRecords[i].NftMint,
			Reverse: v,
		}
	}
	return result, nil
}

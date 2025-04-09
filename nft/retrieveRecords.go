package nft

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

func RetrieveRecords(
	conn *rpc.Client,
	owner solana.PublicKey,
) ([]NftRecord, error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, 2)
	result, err := conn.GetProgramAccountsWithOpts(
		context.TODO(),
		solana.TokenProgramID,
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 32,
						Bytes:  owner.Bytes(),
					},
				},
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 64,
						Bytes:  data,
					},
				},
				{
					DataSize: 165,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("empty result from call to GetProgramAccount")
	}

	fmt.Println("len of results", len(result))

	tokenAcc := make([]token.Account, 0, len(result))
	promises := make([]NftRecord, 0, len(result))

	for _, v := range result {
		var account token.Account
		if err := bin.NewBorshDecoder(v.Account.Data.GetBinary()).
			Decode(&account); err != nil {
			tokenAcc = append(tokenAcc, token.Account{})
			continue
		}
		tokenAcc = append(tokenAcc, account)
	}

	for _, v := range tokenAcc {
		record, err := GetRecordFromMint(conn, v.Mint)
		if err != nil {
			promises = append(promises, NftRecord{})
			continue
		}

		if len(record) == 1 {
			if data := record[0].Data; data != nil {
				nf := NftRecord{}
				if err := borsh.Deserialize(&nf, data.GetBinary()); err != nil {
					promises = append(promises, NftRecord{})
					continue
				}
				promises = append(promises, nf)
				continue
			}
		}
		promises = append(promises, NftRecord{})
	}
	return promises, nil
}

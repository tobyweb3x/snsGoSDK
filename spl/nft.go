package spl

// RetrieveNftOwner returns the publicKey associated with the nameAccount parameter. It returns
// named errors: ErrIgnored and ErrZeroMintSupply, which can be ignored.
/*

ns := someFuncCall()
if nftOwner, err = RetrieveNftOwner(conn, nameAccountKey); err != nil {
		if errors.Is(err, ErrZeroMintSupply) ||
			errors.Is(err, ErrIgnored) {
			return RetrieveResult{
				Registry: ns,
			}, nil
		}
		return RetrieveResult{}, err
	}
*/
// This is due to accounting for the error throwing of Javascript which does not neccesarilly halt program execution.
// func RetrieveNftOwner(conn *rpc.Client, nameAccount solana.PublicKey) (solana.PublicKey, error) {

// 	var (
// 		mint   solana.PublicKey
// 		result rpc.JsonRpcResponse[rpc.GetProgramAccounts]
// 		err    error
// 	)

// 	seeds := [][]byte{
// 		MIntPrefix,
// 		nameAccount.Bytes(),
// 	}

// 	if mint, _, err = solana.FindProgramAddress(seeds, NameTokenizerID); err != nil {
// 		return solana.PublicKey{}, err
// 	}

// 	mintInfo, err := conn.(conn, mint, NoCommitmentArg, NoPublickKeyArg)
// 	if err != nil {
// 		if errors.Is(err, ErrTokenAccountNotFound) || errors.Is(err, ErrTokenInvalidAccountOwner) || errors.Is(err, ErrTokenInvalidAccountSize) {
// 			return mint, ErrIgnored
// 		}
// 		return solana.PublicKey{}, err
// 	}

// 	if mintInfo.Supply == 0 {
// 		return mint, ErrZeroMintSupply
// 	}

// 	filter := rpc.GetProgramAccountsConfig{
// 		Filters: []rpc.GetProgramAccountsConfigFilter{
// 			{
// 				MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
// 					Offset: 0,
// 					Bytes:  mint.String(),
// 				},
// 			},
// 			{
// 				MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
// 					Offset: 64,
// 					Bytes:  "2",
// 				},
// 				// DataSize: 165,
// 			},
// 			{
// 				DataSize: 165,
// 			},
// 		},
// 	}

// 	if result, err = conn.conn.GetProgramAccountsWithConfig(context.Background(), TOKEN_PROGRAM_ID.ToBase58(), filter); err != nil {
// 		return solana.PublicKey{}, err
// 	}

// 	if len(result.GetResult()) != 1 {
// 		return solana.PublicKey{}, fmt.Errorf("unexpected length")
// 	}

// 	if data, ok := result.GetResult()[0].Account.Data.([]byte); ok {
// 		return solana.PublicKeyFromBytes(data[32:64]), nil
// 	}

// 	return solana.PublicKey{}, fmt.Errorf("unexpected data type")
// }

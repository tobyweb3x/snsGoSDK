package twitter_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"snsGoSDK/twitter"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTwitterRegistry(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Fatalf("cannot load env file: error: %s", err.Error())
	}
	conn := rpc.New(os.Getenv("RPC_ENDPOINT"))
	// conn := rpc.New(rpc.MainNetBeta_RPC)
	t.Cleanup(
		func() {
			if err := conn.Close(); err != nil {
				t.Logf("Failed to close connection: %v", err)
			}
		},
	)

	payer := solana.MustPublicKeyFromBase58("JB27XSKgYFBsuxee5yAS2yi1NKSU6WV5GZrKdrzeTHYC")
	tests := []struct {
		name string
		fn   func(*rpc.Client) ([]*solana.GenericInstruction, error)
	}{
		{
			name: "Test case 1/Create & delete instruction",
			fn: func(conn *rpc.Client) ([]*solana.GenericInstruction, error) {
				bytes := make([]byte, 10)
				if _, err := rand.Read(bytes); err != nil {
					return nil, fmt.Errorf("err generating random values: error: %v", err)
				}

				handle := hex.EncodeToString(bytes)
				user := solana.NewWallet().PublicKey()

				ixnsOne, err := twitter.CreateVerifiedTwitterRegistry(
					conn,
					handle,
					user,
					payer,
					10,
				)
				if err != nil {
					return nil, err
				}

				ixnsTwo, err := twitter.DeleteTwitterRegistry(
					handle,
					user,
				)
				if err != nil {
					return nil, err
				}

				ixnSlice := make([]*solana.GenericInstruction, 0, len(ixnsOne)+len(ixnsTwo))
				ixnSlice = append(ixnSlice, ixnsOne...)
				ixnSlice = append(ixnSlice, ixnsTwo...)

				return ixnSlice, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ixns, err := tt.fn(conn)
			if err != nil {
				t.Fatalf("testFn failed: error: %v", err)
				return
			}

			instructions := make([]solana.Instruction, len(ixns))
			for i, ixn := range ixns {
				instructions[i] = ixn
			}

			recent, err := conn.GetLatestBlockhash(context.TODO(), rpc.CommitmentConfirmed)
			if err != nil {
				t.Fatalf("getLatestBlockhash failed: error: %v", err)
				return
			}

			tx, err := solana.NewTransaction(
				instructions,
				recent.Value.Blockhash,
				solana.TransactionPayer(payer),
			)
			if err != nil {
				t.Fatalf("newTransaction failed: error: %v", err)
				return
			}

			_, err = tx.Sign(
				func(key solana.PublicKey) *solana.PrivateKey {
					p := solana.MustPrivateKeyFromBase58(os.Getenv("TEST_PRIVATE_KEY"))
					return &p
				},
			)
			if err != nil {
				t.Fatalf("err signing txn: error: %v", err)
				return
			}

			simTxn, err := conn.SimulateTransactionWithOpts(
				context.TODO(),
				tx,
				&rpc.SimulateTransactionOpts{
					Commitment: rpc.CommitmentConfirmed,
				},
			)
			if err != nil {
				t.Fatalf("simulateTransactionWithOpts failed: error: %v", err)
				return
			}

			assert.Nil(t, simTxn.Value.Err)
			arr, _ := tx.AccountMetaList()
			for _, v := range tx.Message.Instructions {
				fmt.Println("programId---", arr[v.ProgramIDIndex])
				fmt.Println("keys----", len(v.Accounts))
				for _, v := range v.Accounts {
					fmt.Printf("%+v\n", *arr[v])
				}
				fmt.Println("data----", v.Data)
				fmt.Println()
			}

			fmt.Println("Logs:", len(simTxn.Value.Logs))
			for i, v := range simTxn.Value.Logs {
				fmt.Printf("%d--->  %s\n", i+1, v)
			}
		})
	}
}

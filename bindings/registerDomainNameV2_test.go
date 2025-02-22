package bindings_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"slices"
	"snsGoSDK/bindings"
	"snsGoSDK/spl"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestRegisterDomainNameV2(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Fatalf("cannot load env file: error: %s", err.Error())
	}
	conn := rpc.New(os.Getenv("RPC_ENDPOINT"))
	t.Cleanup(
		func() {
			if err := conn.Close(); err != nil {
				t.Logf("Failed to close connection: %v", err)
			}
		},
	)

	tests := []struct {
		Name string
		Fn   func(conn *rpc.Client) ([]*solana.GenericInstruction, error)
	}{
		{
			Name: "Indempotent ATA creation ref",
			Fn: func(conn *rpc.Client) ([]*solana.GenericInstruction, error) {
				buyerTokenAccount, err := spl.GetAssociatedTokenAddressSync(
					spl.PYTHMint,
					spl.VaultOwner,
					true,
				)
				if err != nil {
					return nil, err
				}

				arr := make([]*solana.GenericInstruction, 0, 7)
				for range 3 {
					randomBytes := make([]byte, 10)
					if _, err := rand.Read(randomBytes); err != nil {
						return nil, fmt.Errorf("err generating random values: error: %v", err)
					}

					ixns, err := bindings.RegisterDomainNameV2(
						conn,
						hex.EncodeToString(randomBytes),
						1_000,
						spl.VaultOwner,
						buyerTokenAccount,
						spl.PYTHMint,
						spl.REFERRERS[0],
					)
					if err != nil {
						return nil, err
					}
					arr = append(arr, ixns...)
				}
				return slices.Clip(arr), nil
			},
		},
		{
			Name: "Register V2",
			Fn: func(conn *rpc.Client) ([]*solana.GenericInstruction, error) {
				buyerTokenAccount, err := spl.GetAssociatedTokenAddressSync(
					spl.FIDAMint,
					spl.VaultOwner,
					true,
				)
				if err != nil {
					return nil, err
				}

				arr := make([]*solana.GenericInstruction, 0, 7)

				randomBytes := make([]byte, 10)
				if _, err := rand.Read(randomBytes); err != nil {
					return nil, fmt.Errorf("err generating random values: error: %v", err)
				}

				ixns, err := bindings.RegisterDomainNameV2(
					conn,
					hex.EncodeToString(randomBytes),
					1_000,
					spl.VaultOwner,
					buyerTokenAccount,
					spl.FIDAMint,
					spl.REFERRERS[1],
				)
				if err != nil {
					return nil, err
				}
				arr = append(arr, ixns...)

				return slices.Clip(arr), nil
			},
		}, {
			Name: "Register V2 with ref",
			Fn: func(conn *rpc.Client) ([]*solana.GenericInstruction, error) {
				buyerTokenAccount, err := spl.GetAssociatedTokenAddressSync(
					spl.FIDAMint,
					spl.VaultOwner,
					true,
				)
				if err != nil {
					return nil, err
				}

				arr := make([]*solana.GenericInstruction, 0, 7)

				randomBytes := make([]byte, 10)
				if _, err := rand.Read(randomBytes); err != nil {
					return nil, fmt.Errorf("err generating random values: error: %v", err)
				}

				ixns, err := bindings.RegisterDomainNameV2(
					conn,
					hex.EncodeToString(randomBytes),
					1_000,
					spl.VaultOwner,
					buyerTokenAccount,
					spl.FIDAMint,
					spl.REFERRERS[1],
				)
				if err != nil {
					return nil, err
				}
				arr = append(arr, ixns...)

				return slices.Clip(arr), nil
			},
		},
	}

	for _, tt := range tests {
		// t.Parallel()
		t.Run(tt.Name, func(t *testing.T) {
			ixns, err := tt.Fn(conn)
			if err != nil {
				t.Fatalf("testFn failed: error: %v", err)
			}
			instructions := make([]solana.Instruction, len(ixns))
			for i, ix := range ixns {
				instructions[i] = ix
			}

			recent, err := conn.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
			if err != nil {
				t.Fatalf("getLatestBlockhash failed: error: %v", err)
			}

			tx, err := solana.NewTransaction(
				instructions,
				recent.Value.Blockhash,
				solana.TransactionPayer(spl.VaultOwner),
			)
			if err != nil {
				t.Fatalf("newTransaction failed: error: %v", err)
			}

			_, err = tx.Sign(
				func(key solana.PublicKey) *solana.PrivateKey {
					p := solana.MustPrivateKeyFromBase58(os.Getenv("TEST_PRIVATE_KEY"))
					return &p
				},
			)
			if err != nil {
				t.Fatalf("err signing txn: error: %v", err)
			}

			simTxn, err := conn.SimulateTransactionWithOpts(
				context.TODO(),
				tx,
				&rpc.SimulateTransactionOpts{},
			)
			if err != nil {
				t.Fatalf("simulateTransactionWithOpts failed: error: %v", err)
			}

			assert.Nil(t, simTxn.Value.Err)
			// fmt.Println("Logs:", len(simTxn.Value.Logs))
			// for i, v := range simTxn.Value.Logs {
			// 	fmt.Printf("%d--->  %s\n", i+1, v)
			// }
		})
	}
}

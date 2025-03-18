package bindings_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"snsGoSDK/bindings"
	"snsGoSDK/resolve"
	"snsGoSDK/spl"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubdomain(t *testing.T) {

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
		name string
		fn   func(*rpc.Client) ([]*solana.GenericInstruction, error)
	}{
		{
			name: "resolve & createSubdomain",
			fn: func(conn *rpc.Client) ([]*solana.GenericInstruction, error) {

				sub, parent := "gvbhnjklmjnhb", "bonfida.sol"
				parentOwner, err := resolve.Resolve(conn, parent, resolve.ResolveConfig{})
				if err != nil {
					return nil, fmt.Errorf("resolve() failed: %s", err.Error())
				}
				ix, err := bindings.CreateSubdomain(
					conn,
					fmt.Sprintf("%s.%s", sub, parent),
					1_000,
					parentOwner,
					solana.PublicKey{},
				)
				if err != nil {
					return nil, fmt.Errorf("CreateSubdomain() failed: %s", err.Error())
				}
				return ix, nil
			},
		},
		{
			name: "createsubdomain",
			fn: func(conn *rpc.Client) ([]*solana.GenericInstruction, error) {
				randomBytes := make([]byte, 10)
				if _, err := rand.Read(randomBytes); err != nil {
					return nil, fmt.Errorf("err generating random values: error: %v", err)
				}
				ix, err := bindings.CreateSubdomain(
					conn,
					fmt.Sprintf("%s.bonfida", string(hex.EncodeToString(randomBytes))),
					2_000,
					solana.MustPublicKeyFromBase58("HKKp49qGWXd639QsuH7JiLijfVW5UtCVY4s1n2HANwEA"),
					solana.PublicKey{},
				)
				if err != nil {
					return nil, err
				}
				return ix, nil
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
			for i, ix := range ixns {
				instructions[i] = ix
			}

			recent, err := conn.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
			if err != nil {
				t.Fatalf("getLatestBlockhash failed: error: %v", err)
				return
			}

			tx, err := solana.NewTransaction(
				instructions,
				recent.Value.Blockhash,
				solana.TransactionPayer(spl.VaultOwner),
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

			out, err := conn.SimulateTransactionWithOpts(
				context.TODO(),
				tx,
				&rpc.SimulateTransactionOpts{},
			)
			if err != nil {
				t.Fatalf("simulateTransactionWithOpts failed: error: %v", err)
				return
			}

			assert.Nil(t, out.Value.Err)
		})
	}
}

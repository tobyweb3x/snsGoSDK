package bindings_test

import (
	"context"
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

	testFn := func(conn *rpc.Client) ([]*solana.GenericInstruction, error) {

		sub, parent := "gvbhnjklmjnhb", "bonfida.sol"

		parentOwner, err := resolve.Resolve(conn, parent, resolve.ResolveConfig{})
		if err != nil {
			return nil, fmt.Errorf("resolve() failed in createSubdomain: %s", err.Error())
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
	}

	t.Run("Test case 1", func(t *testing.T) {
		ixns, err := testFn(conn)
		if !assert.Nil(t, err) {
			t.Fatalf("testFn failed: error: %v", err)
		}

		instructions := make([]solana.Instruction, len(ixns))
		for i, ix := range ixns {
			instructions[i] = ix
		}

		recent, err := conn.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if !assert.Nil(t, err) {
			t.Fatalf("getLatestBlockhash failed: error: %v", err)
		}

		tx, err := solana.NewTransaction(
			instructions,
			recent.Value.Blockhash,
			solana.TransactionPayer(spl.VaultOwner),
		)
		if !assert.Nil(t, err) {
			t.Fatalf("newTransaction failed: error: %v", err)
		}

		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				p := solana.MustPrivateKeyFromBase58(os.Getenv("TEST_PRIVATE_KEY"))
				return &p
			},
		)
		if !assert.Nil(t, err) {
			t.Fatalf("err signing txn: error: %v", err)
		}

		out, err := conn.SimulateTransactionWithOpts(
			context.TODO(),
			tx,
			&rpc.SimulateTransactionOpts{},
		)

		if !assert.Nil(t, err) {
			t.Fatalf("simulateTransactionWithOpts failed: error: %v", err)
		}

		assert.Nil(t, out.Value.Err)
	})
}

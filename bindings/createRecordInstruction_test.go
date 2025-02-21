package bindings_test

import (
	"context"
	"os"
	"snsGoSDK/bindings"
	"snsGoSDK/types"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCreateRecordInstruction(t *testing.T) {
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

	owner := solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")
	testFn := func(conn *rpc.Client) (*solana.GenericInstruction, error) {
		return bindings.CreateRecordInstruction(
			conn,
			types.A,
			"wallet-guide-3.sol",
			"192.168.0.1",
			owner,
			owner,
		)
	}

	t.Run("createRecordInstruction", func(t *testing.T) {
		ixn, err := testFn(conn)
		if err != nil {
			t.Fatalf("testFn failed: error: %v", err)
			return
		}

		recent, err := conn.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			t.Fatalf("getLatestBlockhash failed: error: %v", err)
			return
		}

		tx, err := solana.NewTransaction(
			[]solana.Instruction{ixn},
			recent.Value.Blockhash,
			solana.TransactionPayer(owner),
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

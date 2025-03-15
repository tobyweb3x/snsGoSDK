package bindings_test

import (
	"context"
	"os"
	"snsGoSDK/bindings"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestBurnDomain(t *testing.T) {

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

	t.Run("BurnDomain", func(t *testing.T) {

		owner := solana.MustPublicKeyFromBase58("3Wnd5Df69KitZfUoPYZU438eFRNwGHkhLnSAWL65PxJX")
		ix, err := bindings.BurnDomain(
			"1automotive",
			owner,
			owner,
		)
		if err != nil {
			t.Fatal(err)
			return
		}

		recent, err := conn.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			t.Fatalf("getLatestBlockhash failed: error: %v", err)
			return
		}

		tx, err := solana.NewTransactionBuilder().
			AddInstruction(ix).
			SetRecentBlockHash(recent.Value.Blockhash).
			SetFeePayer(owner).Build()
		if err != nil {
			t.Fatal(err)
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
			&rpc.SimulateTransactionOpts{},
		)
		if err != nil {
			t.Fatalf("simulateTransactionWithOpts failed: error: %v", err)
			return
		}

		assert.Nil(t, simTxn.Value.Err)

	})
}

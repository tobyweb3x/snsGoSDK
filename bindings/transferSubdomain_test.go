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

func TestTransferSubdomain(t *testing.T) {
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
		name                string
		feePayer            solana.PublicKey
		isParentOwnerSigner bool
		// fn func(*rpc.Client) ([]*solana.GenericInstruction, error)
	}{
		{
			name:                "Transfer sub - isParentOwnerSigner set to false",
			feePayer:            solana.MustPublicKeyFromBase58("A41TAGFpQkFpJidLwH37ydunE7Q3jpBaS228RkoXiRQk"),
			isParentOwnerSigner: false,
		},
		{
			name:                "Transfer sub - isParentOwnerSigner set to true",
			feePayer:            solana.MustPublicKeyFromBase58("A41TAGFpQkFpJidLwH37ydunE7Q3jpBaS228RkoXiRQk"),
			isParentOwnerSigner: true,
		},
	}

	fn := func(conn *rpc.Client, isParentOwnerSigner bool) ([]*solana.GenericInstruction, error) {
		ixn, err := bindings.TransferSubdomain(
			conn,
			"test.0x33.sol",
			solana.PublicKey{},
			solana.PublicKey{},
			isParentOwnerSigner,
		)
		if err != nil {
			return nil, err
		}
		return []*solana.GenericInstruction{ixn}, nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ixns, err := fn(conn, tt.isParentOwnerSigner)
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
				solana.TransactionPayer(tt.feePayer),
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

			assert.Equal(t, out.Value.Err, "AccountNotFound")
			// assert.Nil(t, out.Value.Err)
		})

	}
}

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

func TestCreateRecordV2Instruction(t *testing.T) {
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
		fn   func() ([]*solana.GenericInstruction, error)
	}{
		{
			name: "Create record V2",
			fn: func() ([]*solana.GenericInstruction, error) {
				owner := solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")
				ix, err := bindings.CreateRecordV2Instruction(
					"wallet-guide-9",
					"bonfida",
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				return []*solana.GenericInstruction{ix}, nil
			},
		},
		{
			name: "Update record V2",
			fn: func() ([]*solana.GenericInstruction, error) {
				domain := "wallet-guide-9"
				owner := solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")
				ixOne, err := bindings.CreateRecordV2Instruction(
					domain,
					"bonfida",
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}

				ixns := make([]*solana.GenericInstruction, 0, 2)
				ixns = append(ixns, ixOne)
				ixTwo, err := bindings.UpdateRecordV2Instruction(
					domain,
					"some text",
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				ixns = append(ixns, ixTwo)
				return ixns, nil
			},
		},
		{
			name: "Delete record V2",
			fn: func() ([]*solana.GenericInstruction, error) {
				domain := "wallet-guide-9"
				owner := solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")
				ixOne, err := bindings.CreateRecordV2Instruction(
					domain,
					"bonfida",
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}

				ixns := make([]*solana.GenericInstruction, 0, 2)
				ixns = append(ixns, ixOne)
				ixTwo, err := bindings.DeleteRecordV2(
					domain,
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				ixns = append(ixns, ixTwo)
				return ixns, nil
			},
		},
		{
			name: "Validate(Solana) record V2",
			fn: func() ([]*solana.GenericInstruction, error) {
				domain := "wallet-guide-9"
				owner := solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")
				ixOne, err := bindings.CreateRecordV2Instruction(
					domain,
					"bonfida",
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}

				ixns := make([]*solana.GenericInstruction, 0, 2)
				ixns = append(ixns, ixOne)
				ixTwo, err := bindings.ValidateRecordV2Content(
					domain,
					true,
					types.Github,
					owner,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				ixns = append(ixns, ixTwo)
				return ixns, nil
			},
		},
		{
			name: "Validate(Eth) record V2",
			fn: func() ([]*solana.GenericInstruction, error) {
				domain := "wallet-guide-9"
				owner := solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")
				ixOne, err := bindings.CreateRecordV2Instruction(
					domain,
					"0x4bfbfd1e018f9f27eeb788160579daf7e2cd7da7",
					types.ETH,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}

				ixns := make([]*solana.GenericInstruction, 0, 3)
				ixns = append(ixns, ixOne)
				ixTwo, err := bindings.ValidateRecordV2Content(
					domain,
					true,
					types.ETH,
					owner,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				ixns = append(ixns, ixTwo)

				ixThree, err := bindings.EthValidateRecordV2Content(
					domain,
					types.ETH,
					owner,
					owner,
					[]byte{78, 235, 200, 2, 51, 5, 225, 127, 83, 156, 25, 226, 53, 239, 196, 189,
						196, 197, 121, 2, 91, 2, 99, 11, 31, 179, 5, 233, 52, 246, 137, 252, 72,
						27, 67, 15, 86, 42, 62, 117, 140, 223, 159, 142, 86, 227, 233, 185, 149,
						111, 92, 122, 147, 23, 217, 1, 66, 72, 63, 150, 27, 219, 152, 10, 28},
					[]byte{75, 251, 253, 30, 1, 143, 159, 39, 238, 183, 136, 22, 5, 121, 218, 247,
						226, 205, 125, 167},
				)
				if err != nil {
					return nil, err
				}
				ixns = append(ixns, ixThree)

				return ixns, nil
			},
		},
		{
			name: "Write ROA record V2",
			fn: func() ([]*solana.GenericInstruction, error) {
				domain := "wallet-guide-9"
				owner := solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")
				ixOne, err := bindings.CreateRecordV2Instruction(
					domain,
					"bonfida",
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}

				ixns := make([]*solana.GenericInstruction, 0, 2)
				ixns = append(ixns, ixOne)
				ixTwo, err := bindings.WritRoaRecordV2(
					domain,
					types.Github,
					owner,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				ixns = append(ixns, ixTwo)
				return ixns, nil
			},
		},
		{
			name: "Create subrecord V2",
			fn: func() ([]*solana.GenericInstruction, error) {
				owner := solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")
				ix, err := bindings.CreateRecordV2Instruction(
					"sub-0.wallet-guide-9",
					"bonfida",
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				return []*solana.GenericInstruction{ix}, nil
			},
		},
		{
			name: "Create record, subrecord, update, validate & delete record V2",
			fn: func() ([]*solana.GenericInstruction, error) {
				domain := "wallet-guide-9"
				owner := solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")
				ixOne, err := bindings.CreateRecordV2Instruction(
					domain,
					"bonfida",
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}

				ixns := make([]*solana.GenericInstruction, 0, 4)
				ixns = append(ixns, ixOne)
				ixTwo, err := bindings.UpdateRecordV2Instruction(
					domain,
					"somethingelse",
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				ixns = append(ixns, ixTwo)
				ixThree, err := bindings.ValidateRecordV2Content(
					domain,
					true,
					types.Github,
					owner,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				ixns = append(ixns, ixThree)
				ixFour, err := bindings.DeleteRecordV2(
					domain,
					types.Github,
					owner,
					owner,
				)
				if err != nil {
					return nil, err
				}
				ixns = append(ixns, ixFour)

				return ixns, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			ixns, err := tt.fn()
			if err != nil {
				t.Fatalf("testFn failed: error: %v", err)
				return
			}

			recent, err := conn.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
			if err != nil {
				t.Fatalf("getLatestBlockhash failed: error: %v", err)
				return
			}

			instructions := make([]solana.Instruction, len(ixns))
			for i, ixn := range ixns {
				instructions[i] = ixn
			}

			tx, err := solana.NewTransaction(
				instructions,
				recent.Value.Blockhash,
				solana.TransactionPayer(solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8")),
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
				&rpc.SimulateTransactionOpts{},
			)
			if err != nil {
				t.Fatalf("simulateTransactionWithOpts failed: error: %v", err)
				return
			}

			assert.Nil(t, simTxn.Value.Err)
		})
	}
}

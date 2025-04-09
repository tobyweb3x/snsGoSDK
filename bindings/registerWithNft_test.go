package bindings_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"snsGoSDK/bindings"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"
	"testing"

	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRegisterWithNFT(t *testing.T) {
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

	t.Run("register with NFT", func(t *testing.T) {
		bytes := make([]byte, 10)
		if _, err := rand.Read(bytes); err != nil {
			t.Fatalf("err generating random values: error: %v", err)
		}
		domain := hex.EncodeToString(bytes)

		out, err := utils.GetDomainKeySync(domain, types.V0)
		if err != nil {
			t.Fatal(err)
		}

		reverse, err := utils.GetReverseKey(domain, false)
		if err != nil {
			t.Fatal(err)
		}

		nftMint := solana.MustPublicKeyFromBase58("7cpq5U6ze5PPcTPVxGifXA8xyDp8rgAJQNwBDj8eWd8w")

		nftMetadata, _, err := solana.FindTokenMetadataAddress(nftMint)
		if err != nil {
			t.Fatal(err)
		}

		masterEdition, _, err := solana.FindProgramAddress(
			[][]byte{
				[]byte("metadata"),
				token_metadata.ProgramID.Bytes(), // Token Metadata Program ID
				nftMint.Bytes(),
				[]byte("edition"),
			},
			token_metadata.ProgramID,
		)
		if err != nil {
			t.Fatal(err)
		}

		// https://solscan.io/collection/3c138f8640f62b62016f8020f0532ff888bb0866363c26fb2241bcf28c0776ad#holders
		holder := solana.MustPublicKeyFromBase58("FiUYY19eXuVcEAHSJ87KEzYjYnfKZm6KbHoVtdQBNGfk")
		source := solana.MustPublicKeyFromBase58("Df9Jz3NrGVd5jjjrXbedwuHbCc1hL131bUXq2143tTfQ")

		ix, err := bindings.RegisterWithNft(
			domain,
			1_000,
			out.PubKey,
			reverse,
			holder,
			source,
			nftMetadata,
			nftMint,
			masterEdition,
		)
		if err != nil {
			t.Fatal(err)
		}

		recent, err := conn.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if !assert.Nil(t, err) {
			t.Fatalf("getLatestBlockhash failed: error: %v", err)
		}

		tx, err := solana.NewTransactionBuilder().
			AddInstruction(ix).
			SetRecentBlockHash(recent.Value.Blockhash).
			SetFeePayer(spl.VaultOwner).Build()
		if err != nil {
			t.Fatal(err)
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

		// fmt.Println(len(simTxn.Value.Logs))
		// for i, v := range simTxn.Value.Logs {
		// 	fmt.Printf("%d --- %s\n", i, v)
		// }
	})
}

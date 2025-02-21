package resolve_test

import (
	"os"
	"snsGoSDK/resolve"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestResolveSolRecordV1(t *testing.T) {
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

	testFn := func(conn *rpc.Client) (solana.PublicKey, error) {
		return resolve.ResolveSolRecordV1(
			conn,
			solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8"),
			"wallet-guide-4",
		)
	}

	t.Run("Check sol record/resolveSolRecordV1", func(t *testing.T) {
		got, err := testFn(conn)
		if err != nil {
			t.Fatalf("testFn failed: error: %v", err)
		}

		assert.Equal(t, "Hf4daCT4tC2Vy9RCe9q8avT68yAsNJ1dQe6xiQqyGuqZ", got.String())
	})
}

package twitter_test

import (
	"os"
	"snsGoSDK/spl"
	"snsGoSDK/twitter"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetTwitterRegistry(t *testing.T) {
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

	t.Run("getTwitterRegistry", func(t *testing.T) {
		got, err :=twitter.GetTwitterRegistry(conn, "plenthor")
		if err != nil {
			t.Fatal(err)
			return
		}

		assert.Equal(t, solana.PublicKey{}.String(), got.Registry.Class.String())
		assert.Equal(t, spl.TwitterRootParentRegistryKey.String(), got.Registry.ParentName.String())
		assert.Equal(t, "JB27XSKgYFBsuxee5yAS2yi1NKSU6WV5GZrKdrzeTHYC", got.Registry.Owner.String())
	})
}

package twitter_test

import (
	"fmt"
	"os"
	"snsGoSDK/twitter"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetHandleAndRegistryKey(t *testing.T) {
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

	type want struct {
		handle   string
		registry solana.PublicKey
	}

	tests := []struct {
		name           string
		verifiedPubKey solana.PublicKey
		want           want
	}{
		{
			name:           "Test case 1",
			verifiedPubKey: solana.MustPublicKeyFromBase58("JB27XSKgYFBsuxee5yAS2yi1NKSU6WV5GZrKdrzeTHYC"),
			want: want{
				handle:   "plenthor",
				registry: solana.MustPublicKeyFromBase58("HrguVp54KnhQcRPaEBULTRhC2PWcyGTQBfwBNVX9SW2i"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("GetHandleAndRegistryKey:%s", tt.name), func(t *testing.T) {
			var (
				w   want
				err error
			)
			if w.registry, w.handle, err = twitter.GetHandleAndRegistryKey(
				conn,
				tt.verifiedPubKey,
			); err != nil {
				t.Fatalf("%s: error: %s\n", tt.name, err.Error())
				return
			}

			assert.Equal(t, tt.want.handle, w.handle)
			assert.Equal(t, tt.want.registry.String(), w.registry.String())
		})
	}
}

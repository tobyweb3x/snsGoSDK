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

func TestReverseTwitterRegistryStateRetrive(t *testing.T) {
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
		name                     string
		reverseTwitterAccountKey solana.PublicKey
		want                     string
	}{
		{
			name:                     "Test case 1",
			reverseTwitterAccountKey: solana.MustPublicKeyFromBase58("C2MB7RDr4wdwSHAPZ8f5qmScYSUHdPKTL6t5meYdcjjW"),
			want:                     "plenthor",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("ReverseTwitterRegistryState Retrive method: %s", tt.name),
			func(t *testing.T) {
				rts := twitter.ReverseTwitterRegistryState{}
				if err := rts.Retrieve(conn, tt.reverseTwitterAccountKey); err != nil {
					t.Fatal(err)
					return
				}

				assert.Equal(t, tt.want, rts.TwitterHandle)
			})
	}
}

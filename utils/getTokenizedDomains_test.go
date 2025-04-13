package utils_test

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"snsGoSDK/utils"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetTokenizedDomains(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Fatalf("cannot load env file: error: %s", err.Error())
	}
	conn := rpc.New(os.Getenv("RPC_ENDPOINT"))
	// conn := rpc.New(rpc.MainNetBeta.RPC)
	t.Cleanup(
		func() {
			if err := conn.Close(); err != nil {
				t.Logf("Failed to close connection: %v", err)
			}
		},
	)

	tests := []struct {
		name  string
		owner solana.PublicKey
		want  []utils.GetTokenizedDomainsResult
	}{
		{
			name:  "Test case 1",
			owner: solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8"),
			want: []utils.GetTokenizedDomainsResult{
				{
					Key:     solana.MustPublicKeyFromBase58("iSNVgWfb31aTWa58UxZ6fp7n3TTrUk5Gojggub5stXk"),
					Mint:    solana.MustPublicKeyFromBase58("2RJhBbxTiPT2bZq5bhjaTZbsnhbDB7VtTAMmCdBrwBZP"),
					Reverse: "wallet-guide-5",
				},
				{
					Key:     solana.MustPublicKeyFromBase58("uDTBDfKrJSBTgmWUZLcENPk5YrHfWbcrUbNFLjsvNpn"),
					Mint:    solana.MustPublicKeyFromBase58("Eskv5Ns4gyREvNPPgANojNPsz6x1cbn9YwT7esAnxPhP"),
					Reverse: "wallet-guide-0",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("getTokenizedDomains:%s", tt.name), func(t *testing.T) {
			got, err := utils.GetTokenizedDomains(conn, tt.owner)
			if err != nil {
				t.Fatalf("getTokenizedDomains failed: error: %s\n", err)
				return
			}
			slices.SortFunc(got, func(a, b utils.GetTokenizedDomainsResult) int {
				return cmp.Compare(b.Reverse, a.Reverse)
			})

			assert.Equal(t, len(tt.want), len(got))
			assert.Equal(t, tt.want, got)
		})
	}
}

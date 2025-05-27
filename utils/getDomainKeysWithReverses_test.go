package utils_test

import (
	"fmt"
	"os"
	"slices"
	"snsGoSDK/utils"
	"strings"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetDomainKeysWithReverses(t *testing.T) {
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
		domainKey,
		domain string
	}

	tests := []struct {
		name string
		user solana.PublicKey
		want []want
	}{
		{
			name: "Test case 1",
			user: solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8"),
			want: []want{
				{domainKey: "9wcWEXmtUbmiAaWdhQ1nSaZ1cmDVdbYNbaeDcKoK5H8r", domain: "wallet-guide-10"},
				{domainKey: "CZFQJkE2uBqdwHH53kBT6UStyfcbCWzh6WHwRRtaLgrm", domain: "wallet-guide-3"},
				{domainKey: "ChkcdTKgyVsrLuD9zkUBoUkZ1GdZjTHEmgh5dhnR4haT", domain: "wallet-guide-4"},
				{domainKey: "2NsGScxHd9bS6gA7tfY3xucCcg6H9qDqLdXLtAYFjCVR", domain: "wallet-guide-6"},
				{domainKey: "6Yi9GyJKoFAv77pny4nxBqYYwFaAZ8dNPZX9HDXw5Ctw", domain: "wallet-guide-7"},
				{domainKey: "8XXesVR1EEsCEePAEyXPL9A4dd9Bayhu9MRkFBpTkibS", domain: "wallet-guide-9"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("GetDomainKeysWithReverses/%s", tt.name), func(t *testing.T) {
			got, err := utils.GetDomainKeysWithReverses(conn, tt.user)
			if err != nil {
				t.Fatalf("GetAllDomain failed: error: %s", err)
				return
			}
			slices.SortFunc(got, func(a, b utils.GetDomainKeysWithReversesResult) int {
				return strings.Compare(a.Domain, b.Domain)
			})
			for i, v := range got {
				assert.Equal(t, tt.want[i].domain, v.Domain)
				assert.Equal(t, tt.want[i].domainKey, v.DomainKey.String())
			}
		})

	}
}

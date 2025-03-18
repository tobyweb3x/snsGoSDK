package utils_test

import (
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

func TestGetAllDomains(t *testing.T) {
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
		user solana.PublicKey
		want []string
	}{
		{
			name: "Test case 1",
			user: solana.MustPublicKeyFromBase58("Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8"),
			want: []string{
				"2NsGScxHd9bS6gA7tfY3xucCcg6H9qDqLdXLtAYFjCVR",
				"6Yi9GyJKoFAv77pny4nxBqYYwFaAZ8dNPZX9HDXw5Ctw",
				"8XXesVR1EEsCEePAEyXPL9A4dd9Bayhu9MRkFBpTkibS",
				"9wcWEXmtUbmiAaWdhQ1nSaZ1cmDVdbYNbaeDcKoK5H8r",
				"CZFQJkE2uBqdwHH53kBT6UStyfcbCWzh6WHwRRtaLgrm",
				"ChkcdTKgyVsrLuD9zkUBoUkZ1GdZjTHEmgh5dhnR4haT",
			},
		},
	}
	for _, tt := range tests {
		// t.Parallel()
		t.Run(fmt.Sprintf("GetAllDomains/%s", tt.name), func(t *testing.T) {
			got, err := utils.GetAllDomains(conn, tt.user)
			if err != nil {
				t.Fatalf("GetAllDomain failed: error: %s", err)
				return
			}
			strSlice := make([]string, len(got))
			for i, v := range got {
				strSlice[i] = v.String()
			}
			slices.Sort(strSlice)
			assert.Equal(t, tt.want, strSlice)
		})

	}

}

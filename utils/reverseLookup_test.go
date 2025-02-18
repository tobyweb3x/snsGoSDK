package utils_test

import (
	"fmt"
	"os"
	"snsGoSDK/utils"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestReverseLookUp(t *testing.T) {
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
		name   string
		domain solana.PublicKey
		want   string
	}{
		{
			name:   "Test case 1",
			domain: solana.MustPublicKeyFromBase58("Crf8hzfthWGbGbLTVCiqRqV5MVnbpHB1L9KQMd6gsinb"),
			want:   "bonfida",
		},
		{
			name:   "Test case 2",
			domain: solana.MustPublicKeyFromBase58("HPjEbJoeS77Qq31tWuS8pZCsY2yHAW2PcpAKBfETuwLa"),
			want:   "tobytobias",
		},
		{
			name:   "Test case 3",
			domain: solana.MustPublicKeyFromBase58("JCqTzrANia2yfS5jDwpM76rFtyVvj4zu2nozVDk29wTh"),
			want:   "menbehindwoman",
		},
		{
			name:   "Test case 4",
			domain: solana.MustPublicKeyFromBase58("2uSQkZRtJDYmBEbSg2WwMeWs2y21PNgGsUNoVRLDGRXZ"),
			want:   "grimmest",
		},
		{
			name:   "Test case 5",
			domain: solana.MustPublicKeyFromBase58("54obixuvJKGeJ6zFwYy1zb55G5c5z3B65MRXcc7fmaVU"),
			want:   "niftydegen",
		},
		// {
		// 	name:   "Test case 6,
		// 	domain: solana.MustPublicKeyFromBase58("5monfqudwcjVztfNa4nAyL4AwwymKgBspdR3RcvhKX4w"),
		// 	want:   "sokka.ambassador",
		// },
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s: - reverseLookUp", tt.name),
			func(t *testing.T) {
				// t.Parallel()
				got, err := utils.ReverseLookup(conn, tt.domain, solana.PublicKey{})
				if assert.Nil(t, err) {
					assert.Equal(t, tt.want, got)
				} else {
					t.Errorf("reverseLookUp error = %v", err)
				}
			})
	}

}

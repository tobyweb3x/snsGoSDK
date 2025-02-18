package utils_test

import (
	"os"
	"snsGoSDK/utils"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestReverseLookUpBatch(t *testing.T) {
	
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

	tests := struct {
		name    string
		domains []solana.PublicKey
		want    []string
	}{
		name: "Test case 1",
		domains: []solana.PublicKey{
			solana.MustPublicKeyFromBase58("Crf8hzfthWGbGbLTVCiqRqV5MVnbpHB1L9KQMd6gsinb"),
			solana.MustPublicKeyFromBase58("HPjEbJoeS77Qq31tWuS8pZCsY2yHAW2PcpAKBfETuwLa"),
			solana.MustPublicKeyFromBase58("JCqTzrANia2yfS5jDwpM76rFtyVvj4zu2nozVDk29wTh"),
			solana.MustPublicKeyFromBase58("2uSQkZRtJDYmBEbSg2WwMeWs2y21PNgGsUNoVRLDGRXZ"),
			solana.MustPublicKeyFromBase58("54obixuvJKGeJ6zFwYy1zb55G5c5z3B65MRXcc7fmaVU"),
			// solana.MustPublicKeyFromBase58("5monfqudwcjVztfNa4nAyL4AwwymKgBspdR3RcvhKX4w"),
		},
		want: []string{
			"bonfida",
			"tobytobias",
			"menbehindwoman",
			"grimmest",
			"niftydegen",
			// "sokka.ambassador",
		},
	}

	t.Run(" reverseLookUpBatch", func(t *testing.T) {
		// t.Parallel()
		got, err := utils.ReverseLookUpBatch(conn, tests.domains)
		if assert.Nil(t, err) {
			assert.Equal(t, tests.want, got)
		} else {
			t.Errorf("reverseLookUpBatch error = %v", err)
		}
	})

}

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

func TestFavoriteDomain(t *testing.T) {

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
		domain  solana.PublicKey
		reverse string
		stale   bool
	}

	tests := []struct {
		name string
		user solana.PublicKey
		want want
	}{
		{
			name: "Test case 1",
			user: solana.MustPublicKeyFromBase58("FidaeBkZkvDqi1GXNEwB8uWmj9Ngx2HXSS5nyGRuVFcZ"),
			want: want{
				domain:  solana.MustPublicKeyFromBase58("Crf8hzfthWGbGbLTVCiqRqV5MVnbpHB1L9KQMd6gsinb"),
				reverse: "bonfida",
				stale:   true,
			},
		},
		{
			name: "Test case 2",
			user: solana.MustPublicKeyFromBase58("HKKp49qGWXd639QsuH7JiLijfVW5UtCVY4s1n2HANwEA"),
			want: want{
				domain:  solana.MustPublicKeyFromBase58("Crf8hzfthWGbGbLTVCiqRqV5MVnbpHB1L9KQMd6gsinb"),
				reverse: "bonfida",
				stale:   false,
			},
		},
		{
			name: "Test case 3",
			user: solana.MustPublicKeyFromBase58("A41TAGFpQkFpJidLwH37ydunE7Q3jpBaS228RkoXiRQk"),
			want: want{
				domain:  solana.MustPublicKeyFromBase58("BaQq8Uib3Aw5SPBedC8MdYCvpfEC9iLkUMHc5M74sAjv"),
				reverse: "1.00728",
				stale:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("GetFavorite domain:%s", tt.name), func(t *testing.T) {
			// t.Parallel()
			got, err := utils.GetFavoriteDoamin(conn, tt.user)
			if err != nil {
				t.Fatal(err)
				return
			}

			assert.Equal(t, tt.want.domain.String(), got.Domain.String())
			assert.Equal(t, tt.want.reverse, got.Reverse)
			assert.Equal(t, tt.want.stale, got.Stale)
		})
	}

}

func TestMultipleFavoriteDomain(t *testing.T) {

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

	test := struct {
		domain []solana.PublicKey
		want   []string
	}{
		domain: []solana.PublicKey{
			// Non tokenized
			solana.MustPublicKeyFromBase58("HKKp49qGWXd639QsuH7JiLijfVW5UtCVY4s1n2HANwEA"),
			// Stale non tokenized
			solana.MustPublicKeyFromBase58("FidaeBkZkvDqi1GXNEwB8uWmj9Ngx2HXSS5nyGRuVFcZ"),
			// Random pubkey
			solana.NewWallet().PublicKey(),
			// Tokenized
			solana.MustPublicKeyFromBase58("36Dn3RWhB8x4c83W6ebQ2C2eH9sh5bQX2nMdkP2cWaA4"),

			solana.MustPublicKeyFromBase58("A41TAGFpQkFpJidLwH37ydunE7Q3jpBaS228RkoXiRQk"),
		},
		want: []string{
			"bonfida",
			"",
			"",
			"fav-tokenized",
			"1.00728",
		},
	}

	t.Run("", func(t *testing.T) {
		// t.Parallel()
		got, err := utils.GetMultipleFavoriteDomain(conn, test.domain)
		if err != nil {
			t.Fatal(err)
			return
		}

		

		fmt.Println(len(got), len(test.want), got)
		assert.True(t, slices.Equal(test.want, got))
	})

}

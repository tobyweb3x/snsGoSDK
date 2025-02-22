package recordv2_test

import (
	"os"
	recordv2 "snsGoSDK/record_v2"
	"snsGoSDK/types"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetMultipleRecordsV2(t *testing.T) {
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

	records := []types.Record{
		types.IPFS,
		types.Email,
		types.Url,
	}

	wants := []string{
		"ipfs://test",
		"test@gmail.com",
		"https://google.com",
	}

	t.Run("GetMultipleRecordsV2", func(t *testing.T) {
		got, err := recordv2.GetMultipleRecordsV2(
			conn,
			"wallet-guide-9.sol",
			records,
			true,
		)
		if err != nil {
			t.Fatalf("GetMultipleRecordsV2 failed: error: %s\n", err.Error())
			return
		}
		ds := make([]string, len(got))
		rs := make([]types.Record, len(got))
		for i, v := range got {
			ds[i] = v.DeserializeContent
			rs[i] = v.Record
		}
		assert.Equal(t, wants, ds)
		assert.Equal(t, records, rs)
	})
}

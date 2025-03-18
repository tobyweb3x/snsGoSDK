package recordv2_test

import (
	"fmt"
	"os"
	recordv2 "snsGoSDK/record_v2"
	"snsGoSDK/types"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetRecordV2(t *testing.T) {
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

	domian := "wallet-guide-9.sol"

	tests := []struct {
		name   string
		record types.Record
		want   string
	}{
		{
			name:   "Test case 1",
			record: types.IPFS,
			want:   "ipfs://test",
		},
		{
			name:   "Test case 2",
			record: types.Email,
			want:   "test@gmail.com",
		},
		{
			name:   "Test case 3",
			record: types.Url,
			want:   "https://google.com",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("getRecordV2:%s", tt.name), func(t *testing.T) {
			got, err := recordv2.GetRecordV2(conn, domian, tt.record, true)
			if err != nil {
				t.Fatalf("getRecordV2:%s failed: error: %s", tt.name, err.Error())
				return
			}
			assert.Equal(t, tt.want, got.DeserializeContent)
		})
	}
}

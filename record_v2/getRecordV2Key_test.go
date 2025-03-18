package recordv2

import (
	"fmt"
	"os"
	"snsGoSDK/types"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetRecordV2Key(t *testing.T) {
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
		domain string
		record types.Record
		want   string
	}{
		{
			name:   "Test case 1",
			domain: "domain1.sol",
			record: types.SOL,
			want:   "GBrd6Q53eu1T2PiaQAtm92r3DwxmoGvZ2D6xjtVtN1Qt",
		},
		{
			name:   "Test case 2",
			domain: "sub.domain2.sol",
			record: types.SOL,
			want:   "A3EFmyCmK5rp73TdgLH8aW49PJ8SJw915arhydRZ6Sws",
		},
		{
			name:   "Test case 3",
			domain: "domain3.sol",
			record: types.Url,
			want:   "DMZmnjcAnUwSje4o2LGJhipCfNZ5b37GEbbkwbQBWEW1",
		},
		{
			name:   "Test case 4",
			domain: "sub.domain4.sol",
			record: types.Url,
			want:   "6o8JQ7vss6r9sw9GWNVugZktwfEJ67iUz6H63hhmg4sj",
		},
		{
			name:   "Test case 5",
			domain: "domain5.sol",
			record: types.IPFS,
			want:   "DQHeVmAj9Nz4uAn2dneEsgBZWcfhUqLdtbDcfWhGL47D",
		},
		{
			name:   "Test case 6",
			domain: "sub.domain6.sol",
			record: types.IPFS,
			want:   "Dj7tnTTaktrrmdtatRuLG3YdtGZk8XEBMb4w5WtCBHvr",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("GetRecordV2Key:%s", tt.name), func(t *testing.T) {
			got, err := GetRecordV2Key(tt.domain, tt.record)
			if err != nil {
				t.Fatalf("test failed: error: %v\n", err)
				return
			}
			assert.Equal(t, tt.want, got.String())
		})
	}
}

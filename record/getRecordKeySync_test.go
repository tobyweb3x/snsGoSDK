package record_test

import (
	"os"
	"snsGoSDK/record"
	"snsGoSDK/types"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetRecordKeySync(t *testing.T) {
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
		name    string
		domain  string
		record  types.Record
		expcted string
	}{
		{
			name:    "GetRecordKeySync/domain1.sol",
			domain:  "domain1.sol",
			record:  types.SOL,
			expcted: "ATH9akc5pi1PWDB39YY7VCoYzCxmz8XVj23oegSoNSPL",
		},
		{
			name:    "GetRecordKeySync/sub.domain2.sol",
			domain:  "sub.domain2.sol",
			record:  types.SOL,
			expcted: "AEgJVf6zaQfkyYPnYu8Y9Vxa1Sy69EtRSP8iGubx5MnC",
		},
		{
			name:    "GetRecordKeySync/domain3.sol",
			domain:  "domain3.sol",
			record:  types.Url,
			expcted: "EuxtWLCKsdpwM8ftKjnD2Q8vBdzZunh7DY1mHwXhLTqx",
		},
		{
			name:    "GetRecordKeySync/sub.domain4.sol",
			domain:  "sub.domain4.sol",
			record:  types.Url,
			expcted: "64nv6HSbifdUgdWst48V4YUB3Y3uQXVQRD4iDZPd9qGx",
		},
		{
			name:    "GetRecordKeySync/domain5.sol",
			domain:  "domain5.sol",
			record:  types.IPFS,
			expcted: "2uRMeYzKXaYgFVQ1Yh7fKyZWcxsFUMgpEwMi19sVjwjk",
		},
		{
			name:    "GetRecordKeySync/sub.domain6.sol",
			domain:  "sub.domain6.sol",
			record:  types.IPFS,
			expcted: "61JdnEhbd2bEfxnu2uQ38gM2SUry2yY8kBMEseYh8dDy",
		},
	}

	for _, vv := range tests {
		t.Run(vv.name, func(t *testing.T) {
			got, err := record.GetRecordKeySync(vv.domain, vv.record)
			if err != nil {
				t.Fatalf("%s failed: error: %v", vv.name, err)
				return
			}

			assert.Equal(t, vv.expcted, got.String())
		})
	}
}

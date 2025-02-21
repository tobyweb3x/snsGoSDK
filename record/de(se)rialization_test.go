package record_test

import (
	"os"
	"snsGoSDK/record"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestDeserializeSerializeRcord(t *testing.T) {
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
		content string
		record  types.Record
	}{
		{
			name:    "de(se)rialization for TXT",
			content: "this is a test",
			record:  types.TXT,
		},
		{
			name:    "de(se)rialization for Injective",
			content: "inj13glcnaum2xqv5a0n0hdsmv0f6nfacjsfvrh5j9",
			record:  types.Injective,
		},
		{
			name:    "de(se)rialization for CNAME",
			content: "example.com",
			record:  types.CNAME,
		},
		{
			name:    "de(se)rialization for CNAME",
			content: "example.com",
			record:  types.CNAME,
		},
		{
			name:    "de(se)rialization for ETH",
			content: "0xc0ffee254729296a45a3885639ac7e10f9d54979",
			record:  types.ETH,
		},
		{
			name:    "de(se)rialization for A",
			content: "1.1.1.4",
			record:  types.A,
		},
		{
			name:    "de(se)rialization for AAAA",
			content: "2345:425:2ca1::567:5673:23b5",
			record:  types.AAAA,
		},
		{
			name:    "de(se)rialization for Discord",
			content: "username",
			record:  types.Discord,
		},
	}

	for _, vv := range tests {
		t.Run(vv.name, func(t *testing.T) {
			byteSlice, err := record.SerializeRecord(vv.content, vv.record)
			if err != nil {
				t.Fatalf("serialization err: %s", err.Error())
				return
			}
			nm := spl.NameRegistryState{
				Data: byteSlice,
			}
			str, err := record.DeserializeRecord(nm, vv.record, solana.PublicKey{})
			if err != nil {
				t.Fatalf("deserialization err: %s", err.Error())
				return
			}

			assert.Equal(t, vv.content, str)
		})
	}
}

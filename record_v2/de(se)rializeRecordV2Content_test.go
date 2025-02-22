package recordv2_test

import (
	"os"
	recordv2 "snsGoSDK/record_v2"
	"snsGoSDK/types"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestDeserializeSerializeRecordV2Content(t *testing.T) {
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
		length  int
	}{
		{
			name:    "de(se)rializationRecordV2Content for TXT",
			content: "this is test",
			record:  types.TXT,
		},
		{
			name:    "de(se)rializationRecordV2Content for SOL",
			content: solana.NewWallet().PublicKey().String(),
			record:  types.SOL,
			length:  32,
		},
		{
			name:    "de(se)rializationRecordV2Content for Injective",
			content: "inj13glcnaum2xqv5a0n0hdsmv0f6nfacjsfvrh5j9",
			record:  types.Injective,
			length:  20,
		},
		{
			name:    "de(se)rializationRecordV2Content for CNAME",
			content: "example.com",
			record:  types.CNAME,
		},
		{
			name:    "de(se)rializationRecordV2Content for CNAME",
			content: "example.com",
			record:  types.CNAME,
		},
		{
			name:    "de(se)rializationRecordV2Content for ETH",
			content: "0xc0ffee254729296a45a3885639ac7e10f9d54979",
			record:  types.ETH,
			length:  20,
		},
		{
			name:    "de(se)rializationRecordV2Content for A",
			content: "1.1.1.4",
			record:  types.A,
			length:  4,
		},
		{
			name:    "de(se)rializationRecordV2Content for AAAA",
			content: "2345:425:2ca1::567:5673:23b5",
			record:  types.AAAA,
			length:  16,
		},
		{
			name:    "de(se)rializationRecordV2Content for Discord",
			content: "username",
			record:  types.Discord,
		},
		{
			name:    "de(se)rializationRecordV2Content for IPNS",
			content: "k51qzi5uqu5dlvj2baxnqndepeb86cbk3ng7n3i46uzyxzyqj2xjonzllnv0v8",
			record:  types.IPNS,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ser, err := recordv2.SerializeRecordV2Content(test.content, test.record)
			if err != nil {
				t.Fatalf("ser. failed: error: %s\n", err.Error())
				return
			}

			des, err := recordv2.DeserializeRecordV2Content(ser, test.record)
			if err != nil {
				t.Fatalf("des. failed: error: %s\n", err.Error())
				return
			}

			assert.Equal(t, test.content, des)
			if test.length > 0 {
				assert.Equal(t, test.length, len(ser))
			}
		})
	}
}

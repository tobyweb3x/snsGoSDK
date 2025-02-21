package record_test

import (
	"fmt"
	"os"
	"snsGoSDK/record"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRecords(t *testing.T) {

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
		expect string
		fn     func(conn *rpc.Client, domain string) (string, error)
	}{
		{
			name:   "IPFS",
			domain: "🍍",
			expect: "QmbWqxBEKC3P8tqsKc98xmWNzrzDtRLMiMPL8wBuTGsMnR",
			fn:     record.GetIPFSRecord,
		},
		{
			name:   "Arweave",
			domain: "🍍",
			expect: "some-arweave-hash",
			fn:     record.GetArweaveRecord,
		},
		{
			name:   "Ethereum",
			domain: "🍍",
			expect: "0x570eDC13f9D406a2b4E6477Ddf75D5E9cCF51cd6",
			fn:     record.GetETHRecord,
		},
		{
			name:   "Bitcoin",
			domain: "🍍",
			expect: "3JfBcjv7TbYN9yQsyfcNeHGLcRjgoHhV3z",
			fn:     record.GetBTCRecord,
		},
		{
			name:   "Litecoin",
			domain: "🍍",
			expect: "MK6deR3Mi6dUsim9M3GPDG2xfSeSAgSrpQ",
			fn:     record.GetLTCRecord,
		},
		{
			name:   "Dogecoin",
			domain: "🍍",
			expect: "DC79kjg58VfDZeMj9cWNqGuDfYfGJg9DjZ",
			fn:     record.GetDogeRecord,
		},
		{
			name:   "Email",
			domain: "🍍",
			expect: "🍍@gmail.com",
			fn:     record.GetEmailRecord,
		},
		{
			name:   "URL",
			domain: "🍍",
			expect: "🍍.io",
			fn:     record.GetURLRecord,
		},
		{
			name:   "Discord",
			domain: "🍍",
			expect: "@🍍#7493",
			fn:     record.GetDiscordRecord,
		},
		{
			name:   "GitHub",
			domain: "🍍",
			expect: "@🍍_dev",
			fn:     record.GetGithubRecord,
		},
		{
			name:   "Reddit",
			domain: "🍍",
			expect: "@reddit-🍍",
			fn:     record.GetRedditRecord,
		},
		{
			name:   "Twitter",
			domain: "🍍",
			expect: "@🍍",
			fn:     record.GetTwitterRecord,
		},
		{
			name:   "Telegram",
			domain: "🍍",
			expect: "@🍍-tg",
			fn:     record.GetTelegramRecord,
		},
		{
			name:   "BSC",
			domain: "aanda.sol",
			expect: "0x4170ad697176fe6d660763f6e4dfcf25018e8b63",
			fn:     record.GetBSCRecord,
		},
		{
			name:   "subdomain with emailRecord",
			domain: "test.🇺🇸.sol",
			expect: "test@test.com",
			fn:     record.GetEmailRecord,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Records:%s", test.name),
			func(t *testing.T) {
				// t.Parallel()
				got, err := test.fn(conn, test.domain)
				if err != nil {
					t.Fatalf("failed to get record-%s: error: %v", test.name, err)
					return
				}

				assert.Equal(t, test.expect, got)
			})
	}
}

package resolve_test

import (
	"fmt"
	"os"
	"testing"

	"snsGoSDK/resolve"
	"snsGoSDK/spl"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
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
		config resolve.ResolveConfig
		want   string
	}{
		{
			name:   "Test case 1",
			domain: "sns-ip-5-wallet-1",
			config: resolve.ResolveConfig{},
			want:   "ALd1XSrQMCPSRayYUoUZnp6KcP6gERfJhWzkP49CkXKs",
		},
		{
			name:   "Test case 2",
			domain: "sns-ip-5-wallet-2",
			config: resolve.ResolveConfig{},
			want:   "AxwzQXhZNJb9zLyiHUQA12L2GL7CxvUNrp6neee6r3cA",
		},
		{
			name:   "Test case 3",
			domain: "sns-ip-5-wallet-4",
			config: resolve.ResolveConfig{},
			want:   "7PLHHJawDoa4PGJUK3mUnusV7SEVwZwEyV5csVzm86J4",
		},
		{
			name:   "Test case 4",
			domain: "sns-ip-5-wallet-5",
			config: resolve.ResolveConfig{AllowPda: resolve.AllowPDATrue, ProgramIDs: []solana.PublicKey{solana.SystemProgramID}},
			want:   "96GKJgm2W3P8Bae78brPrJf4Yi9AN1wtPJwg2XVQ2rMr",
		},
		{
			name:   "Test case 5",
			domain: "sns-ip-5-wallet-5",
			config: resolve.ResolveConfig{AllowPda: resolve.AllowPDAAny},
			want:   "96GKJgm2W3P8Bae78brPrJf4Yi9AN1wtPJwg2XVQ2rMr",
		},
		{
			name:   "Test case 6",
			domain: "sns-ip-5-wallet-7",
			config: resolve.ResolveConfig{},
			want:   "53Ujp7go6CETvC7LTyxBuyopp5ivjKt6VSfixLm1pQrH",
		},
		{
			name:   "Test case 7",
			domain: "sns-ip-5-wallet-8",
			config: resolve.ResolveConfig{},
			want:   "ALd1XSrQMCPSRayYUoUZnp6KcP6gERfJhWzkP49CkXKs",
		},
		{
			name:   "Test case 8",
			domain: "sns-ip-5-wallet-9",
			config: resolve.ResolveConfig{},
			want:   "ALd1XSrQMCPSRayYUoUZnp6KcP6gERfJhWzkP49CkXKs",
		},
		{
			name:   "Test case 9",
			domain: "sns-ip-5-wallet-10",
			config: resolve.ResolveConfig{AllowPda: resolve.AllowPDATrue, ProgramIDs: []solana.PublicKey{solana.SystemProgramID}},
			want:   "96GKJgm2W3P8Bae78brPrJf4Yi9AN1wtPJwg2XVQ2rMr",
		},
		{
			name:   "Test case 10",
			domain: "sns-ip-5-wallet-10",
			config: resolve.ResolveConfig{AllowPda: resolve.AllowPDAAny},
			want:   "96GKJgm2W3P8Bae78brPrJf4Yi9AN1wtPJwg2XVQ2rMr",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s: - resolving domain", tt.name),
			func(t *testing.T) {
				// t.Parallel()
				got, err := resolve.Resolve(conn, tt.domain, tt.config)
				if err != nil {
					t.Fatalf("Resolve() error = %v", err)
					return
				}

				assert.Equal(t, tt.want, got.String())
			})
	}
	_ = tests

	errorTests := []struct {
		name   string
		domain string
		err    string
	}{
		{
			name:   "Test case 1",
			domain: "sns-ip-5-wallet-3",
			err:    string(spl.WrongValidation),
		},
		{
			name:   "Test case 2",
			domain: "sns-ip-5-wallet-6",
			err:    string(spl.PdaOwnerNotAllowed),
		},
		{
			name:   "Test case 3",
			domain: "sns-ip-5-wallet-11",
			err:    string(spl.PdaOwnerNotAllowed),
		},
		{
			name:   "Test case 4",
			domain: "sns-ip-5-wallet-12",
			err:    string(spl.WrongValidation),
		},
	}

	for _, tt := range errorTests {
		t.Run(fmt.Sprintf("%s - should throw named error", tt.name),
			func(t *testing.T) {
				// t.Parallel()
				_, err := resolve.Resolve(conn, tt.domain, resolve.ResolveConfig{})
				fmt.Println(err.Error())
				assert.ErrorContains(t, err, tt.err)
			})
	}
	_ = errorTests

	additionalTests := []struct {
		name   string
		domain string
		owner  string
	}{
		{
			name:   "Test case 1",
			domain: "wallet-guide-5.sol",
			owner:  "Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8",
		},
		{
			name:   "Test case 2",
			domain: "wallet-guide-4.sol",
			owner:  "Hf4daCT4tC2Vy9RCe9q8avT68yAsNJ1dQe6xiQqyGuqZ",
		},
		{
			name:   "Test case 3",
			domain: "wallet-guide-3.sol",
			owner:  "Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8",
		},
		{
			name:   "Test case 4",
			domain: "wallet-guide-2.sol",
			owner:  "36Dn3RWhB8x4c83W6ebQ2C2eH9sh5bQX2nMdkP2cWaA4",
		},
		{
			name:   "Test case 5",
			domain: "wallet-guide-1.sol",
			owner:  "36Dn3RWhB8x4c83W6ebQ2C2eH9sh5bQX2nMdkP2cWaA4",
		},
		{
			name:   "Test case 6",
			domain: "wallet-guide-0.sol",
			owner:  "Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8",
		},
		{
			name:   "Test case 7",
			domain: "sub-0.wallet-guide-3.sol",
			owner:  "Fxuoy3gFjfJALhwkRcuKjRdechcgffUApeYAfMWck6w8",
		},
		{
			name:   "Test case 8",
			domain: "sub-1.wallet-guide-3.sol",
			owner:  "Hf4daCT4tC2Vy9RCe9q8avT68yAsNJ1dQe6xiQqyGuqZ",
		},
		// Record V2
		{
			name:   "Test case 9",
			domain: "wallet-guide-6",
			owner:  "Hf4daCT4tC2Vy9RCe9q8avT68yAsNJ1dQe6xiQqyGuqZ",
		},
		{
			name:   "Test case 10",
			domain: "wallet-guide-8",
			owner:  "36Dn3RWhB8x4c83W6ebQ2C2eH9sh5bQX2nMdkP2cWaA4",
		},
	}

	for _, tt := range additionalTests {
		t.Run(fmt.Sprintf("%s - (checking for backward compatibility)", tt.name),
			func(t *testing.T) {
				// t.Parallel()
				got, err := resolve.Resolve(conn, tt.domain, resolve.ResolveConfig{})
				if assert.Nil(t, err) {
					assert.Equal(t, tt.owner, got.String())
				} else {
					t.Errorf("Resolve() error = %v", err)
				}
			})
	}
	_ = additionalTests
}

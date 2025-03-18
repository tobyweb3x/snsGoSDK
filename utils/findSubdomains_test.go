package utils_test

import (
	"fmt"
	"os"
	"snsGoSDK/types"
	"snsGoSDK/utils"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestFindSubdomains(t *testing.T) {
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
		want   []string
	}{
		{
			name:   "Test case 1",
			domain: "bonfida",
			want:   []string{"naming", "test"},
		},
	}

	fn := func(conn *rpc.Client, domain string) ([]string, error) {
		out, err := utils.GetDomainKeySync(domain, types.VersionUnspecified)
		if err != nil {
			return nil, fmt.Errorf("getDomainKeySync err: %s", err.Error())
		}

		return utils.FindSubdomains(conn, out.PubKey)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fn(conn, tt.domain)
			if err != nil {
				t.Fatal(err)
				return
			}

			assert.Equal(t, tt.want, got)

		})
	}
}

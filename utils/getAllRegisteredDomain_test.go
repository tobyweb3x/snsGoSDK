package utils_test

import (
	"fmt"
	"os"
	"snsGoSDK/utils"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetAllRegisteredDomain(t *testing.T) {

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

	t.Run("getAllRegisteredDomain", func(t *testing.T) {
		got, err := utils.GetAllRegisteredDomain(conn)
		if err != nil {
			t.Fatal(err)
			return
		}
		fmt.Println(len(got))
		assert.True(t, len(got) > 130_000)
	})
}

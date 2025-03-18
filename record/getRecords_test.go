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

func TestGetRecords(t *testing.T) {
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

	t.Run("Get Multiple Records", func(t *testing.T) {
		tests := struct {
			records []types.Record
			expect  []string
		}{
			records: []types.Record{types.Telegram, types.Github, types.Backpack},
			expect:  []string{"@🍍-tg", "@🍍_dev", ""},
		}

		got, err := record.GetRecordsDeserialized(
			conn,
			"🍍",
			tests.records,
		)
		if err != nil {
			t.Fatalf("GetRecordsDeserialized() failed: error: %v", err)
			return
		}
		assert.Equal(t, tests.expect, got)
	})

}

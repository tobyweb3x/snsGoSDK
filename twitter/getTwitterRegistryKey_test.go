package twitter_test

import (
	"snsGoSDK/twitter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTwitterRegistryKey(t *testing.T){

	tests := []struct {
		handle string
		want   string
	}{
		{
			handle: "plenthor",
			want:   "HrguVp54KnhQcRPaEBULTRhC2PWcyGTQBfwBNVX9SW2i",
		},
	}

	for _, tt := range tests {
		// t.Parallel()
		t.Run("GetTwitterRegistryKey", func(t *testing.T) {
			got, err := twitter.GetTwitterRegistryKey(tt.handle)
			if err != nil {
				t.Fatal(err)
				return
			}

			assert.Equal(t, tt.want, got.String())
		})
	}
}


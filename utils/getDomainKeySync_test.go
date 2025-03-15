package utils_test

import (
	"fmt"
	"snsGoSDK/types"
	"snsGoSDK/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDomainKeySync(t *testing.T) {

	tests := []struct {
		name   string
		domain string
		want   string
	}{
		{
			name:   "Test case 1",
			domain: "bonfida",
			want:   "Crf8hzfthWGbGbLTVCiqRqV5MVnbpHB1L9KQMd6gsinb",
		},
		{
			name:   "Test case 2",
			domain: "bonfida.sol",
			want:   "Crf8hzfthWGbGbLTVCiqRqV5MVnbpHB1L9KQMd6gsinb",
		},
		{
			name:   "Test case 3",
			domain: "tobytobias.sol",
			want:   "HPjEbJoeS77Qq31tWuS8pZCsY2yHAW2PcpAKBfETuwLa",
		},
		{
			name:   "Test case 4",
			domain: "dex.bonfida",
			want:   "HoFfFXqFHAC8RP3duuQNzag1ieUwJRBv1HtRNiWFq4Qu",
		},
		{
			name:   "Test case 5",
			domain: "dex.bonfida.sol",
			want:   "HoFfFXqFHAC8RP3duuQNzag1ieUwJRBv1HtRNiWFq4Qu",
		},
		{
			name:   "Test case 6",
			domain: "sub-0.wallet-guide-3.sol",
			want:   "B6qDLYop4KAGbx7JYN41chDoYjKj3Nqc4sUTmbbiTW4v",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("GetDomainKeySync:%s", tt.name), func(t *testing.T) {
			got, err := utils.GetDomainKeySync(tt.domain, types.VersionUnspecified)
			if err != nil {
				t.Fatal(err)
				return
			}

			assert.Equal(t, tt.want, got.PubKey.String())
		})
	}

}

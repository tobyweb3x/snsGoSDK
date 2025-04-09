package nft_test

import (
	"snsGoSDK/nft"
	"snsGoSDK/types"
	"snsGoSDK/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDomainMint(t *testing.T) {

	tests := []struct {
		domain string
		want   string
	}{
		{
			domain: "domain1.sol",
			want:   "3YTxXhhVue9BVjgjPwJbbJ4uGPsnwN453DDf72rYE5WN",
		},
		{
			domain: "sub.domain2.sol",
			want:   "66CnogoXDBqYeYRGYzQf19VyrMnB4uGxpZQDuDYfbKCX",
		},
	}

	for _, tt := range tests {
		domain, err := utils.GetDomainKeySync(tt.domain, types.V0)
		if err != nil {
			t.Fatal(err)
			return
		}

		got, _, err := nft.GetDomainMint(domain.PubKey)
		if err != nil {
			t.Fatal(err)
			return
		}

		assert.Equal(t, tt.want, got.String())
	}
}

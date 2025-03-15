package utils_test

import (
	"fmt"
	"snsGoSDK/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDomainPriceFromName(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{
			name: "Test case 1",
			text: "1",
			want: 750,
		},
		{
			name: "Test case 2",
			text: "✅",
			want: 750,
		},
		{
			name: "Test case 3",
			text: "요",
			want: 750,
		},
		{
			name: "Test case 4",
			text: "👩‍👩‍👧",
			want: 750,
		},
		{
			name: "Test case 5",
			text: "10",
			want: 700,
		},
		{
			name: "Test case 6",
			text: "1✅",
			want: 700,
		},
		{
			name: "Test case 7",
			text: "👩‍👩‍👧✅",
			want: 700,
		},
		{
			name: "Test case 8",
			text: "독도",
			want: 700,
		},
		{
			name: "Test case 9",
			text: "100",
			want: 640,
		},
		{
			name: "Test case 10",
			text: "10✅",
			want: 640,
		},
		{
			name: "Test case 11",
			text: "1독도",
			want: 640,
		},
		{
			name: "Test case 12",
			text: "1000",
			want: 160,
		},
		{
			name: "Test case 13",
			text: "100✅",
			want: 160,
		},
		{
			name: "Test case 14",
			text: "10000",
			want: 20,
		},
		{
			name: "Test case 15",
			text: "1000✅",
			want: 20,
		},
		{
			name: "Test case 16",
			text: "fêtes",
			want: 20,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("GetDomainPriceFromName/%s", tt.name), func(t *testing.T) {
			got := utils.GetDomainPriceFromName(tt.text)
			assert.Equal(t, tt.want, got)
		})
	}
}

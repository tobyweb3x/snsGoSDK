package utils

import (
	"github.com/rivo/uniseg"
)

func GetDomainPriceFromName(name string) int {
	graphemes := uniseg.NewGraphemes(name)

	length := 0
	for graphemes.Next() {
		length++
	}
	switch length {
	case 1:
		return 750
	case 2:
		return 700
	case 3:
		return 640
	case 4:
		return 160
	default:
		return 20
	}
}

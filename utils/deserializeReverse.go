package utils

import (
	"encoding/binary"
	"errors"
	"strings"
)

func DeserializeReverse(data []byte, trimFirstNullByte bool) (string, error) {

	if len(data) == 0 {
		return "", errors.New("data is empty")
	}

	if len(data) < 4 {
		return "", errors.New("data length is less than expected (4)")
	}

	nameLength := binary.LittleEndian.Uint32(data[:4])

	if int(nameLength) > len(data[4:]) {
		return "", errors.New("unexpected data length")
	}

	nameStr := string(data[4 : 4+nameLength])

	if trimFirstNullByte && len(nameStr) > 0 && strings.HasPrefix(nameStr, "\x00") {
		nameStr = strings.TrimPrefix(nameStr, "\x00")
	}

	return nameStr, nil
}

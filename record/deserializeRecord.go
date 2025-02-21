package record

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"slices"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"golang.org/x/net/idna"
)

func DeserializeRecord(
	registry spl.NameRegistryState,
	record types.Record,
	recordKey solana.PublicKey,
) (string, error) {

	buffer := registry.Data
	if len(buffer) == 0 {
		return "", errors.New("empty registry.Data field")
	}
	if slices.Compare(buffer, make([]byte, len(buffer))) == 0 {
		return "", errors.New("empty registry.data")
	}

	size, ok := types.RecordV1Size[record]
	idx := trimNullPaddingIdx(buffer)

	if !ok {
		str := string(buffer[:idx])
		if record == types.CNAME || record == types.TXT {
			str, err := idna.ToASCII(str)
			if err != nil {
				return "", err
			}
			return strings.TrimPrefix(str, "xn--"), nil
		}

		return str, nil
	}

	// Handle SOL record first whether it's over allocated or not
	if record == types.SOL {
		var expectedBuffer bytes.Buffer
		expectedBuffer.Write(buffer[0:32])
		expectedBuffer.Write(recordKey.Bytes())
		expectedHex := hex.EncodeToString(expectedBuffer.Bytes())
		expected := []byte(expectedHex)
		valid := CheckSolRecord(
			expected,
			buffer[32:96],
			registry.Owner,
		)
		if valid {
			return base58.Encode(buffer[0:32]), nil
		}
	}
	fmt.Println("idx", idx, size)
	// Old record UTF-8 encoded
	if ok && int(size) != idx {
		if address := string(buffer[0:idx]); record == types.Injective {
			hrp, decoded, err := bech32.Decode(address)
			if err != nil {
				return "", err
			}
			if strings.HasPrefix(hrp, "inj") && len(decoded) == 20 {
				return address, nil
			}
		} else if record == types.BSC || record == types.ETH {
			hex, err := hex.DecodeString(address[2:])
			if err != nil {
				return "", err
			}
			if strings.HasPrefix(address, "0x") && len(hex) == 20 {
				return address, nil
			}
		} else if record == types.A || record == types.AAAA {
			if ip := net.ParseIP(address); ip != nil {
				return address, nil
			}
		}

		return "", spl.NewSNSError(spl.InvalidRecordData, "the record is malformed", nil)
	}

	if record == types.ETH || record == types.BSC {
		return fmt.Sprintf("0x%s", hex.EncodeToString(buffer[0:size])), nil
	} else if record == types.Injective {
		words, err := bech32.ConvertBits(buffer[0:size], 8, 5, true)
		if err != nil {
			return "", err
		}
		return bech32.Encode("inj", words)
	} else if record == types.A || record == types.AAAA {
		ip := net.IP(buffer[0:size])
		if ip.To4() == nil && ip.To16() == nil {
			return "", spl.NewSNSError(spl.InvalidRecordData, "Invalid IP address", nil)
		}
		return ip.String(), nil
	} else if record == types.Background {
		return solana.PublicKeyFromBytes(buffer[0:size]).String(), nil
	}

	return "", spl.NewSNSError(spl.InvalidRecordData, "the record is malformed", nil)

}

func trimNullPaddingIdx(buffer []byte) int {
	reversed := make([]byte, len(buffer))
	copy(reversed, buffer)
	slices.Reverse(reversed)
	lastNonNull := len(reversed) - 1 - slices.IndexFunc(
		reversed, func(s byte) bool {
			return s != 0
		},
	)
	_ = reversed
	return lastNonNull + 1
}

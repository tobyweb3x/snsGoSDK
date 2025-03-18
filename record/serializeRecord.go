package record

import (
	"encoding/hex"
	"net"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/gagliardetto/solana-go"
	"golang.org/x/net/idna"
)

func SerializeRecord(str string, record types.Record) ([]byte, error) {
	var err error

	if _, ok := types.RecordV1Size[record]; !ok {
		if record == types.CNAME || record == types.TXT {
			if str, err = idna.ToASCII(str); err != nil {
				return nil, nil
			}
			str = strings.TrimPrefix(str, "xn--")
		}

		return []byte(str), nil
	}

	if record == types.SOL {
		return nil, spl.NewSNSError(spl.UnsupportedRecord, "Use `serializeSolRecord` for SOL record", nil)

	} else if record == types.ETH || record == types.BSC {
		if !strings.HasPrefix(str, "0x") {
			return nil, spl.NewSNSError(spl.InvalidEvmAddress, "the record content must start with `0x`", nil)
		}
		str = strings.TrimPrefix(str, "0x")
		return hex.DecodeString(str)

	} else if record == types.Injective {
		hrp, decoded, err := bech32.Decode(str)
		if err != nil {
			return nil, spl.NewSNSError(spl.InvalidInjectiveAddress, "the record content must be a valid bech32 string", err)
		}
		if hrp != "inj" {
			return nil, spl.NewSNSError(spl.InvalidInjectiveAddress, "the record content must start with `inj`", nil)
		}

		if decoded, err = bech32.ConvertBits(decoded, 5, 8, false); err != nil {
			return nil, spl.NewSNSError(spl.InvalidInjectiveAddress, "error coverting bech32 5-bit to 8-bit string", err)
		}

		if len(decoded) != 20 {
			return nil, spl.NewSNSError(spl.InvalidInjectiveAddress, "the record content must be 20 bytes long", nil)
		}
		return decoded, nil

	} else if record == types.A {
		ip := net.ParseIP(str)
		if ip == nil {
			return nil, spl.NewSNSError(spl.InvalidARecord, "the record content must be a valid IP address", nil)
		}

		ip = ip.To4()
		if ip == nil || len(ip) != 4 {
			return nil, spl.NewSNSError(spl.InvalidARecord, "the record content must be a valid IPv4 address", nil)
		}
		return ip, nil

	} else if record == types.AAAA {
		ip := net.ParseIP(str)
		if ip == nil {
			return nil, spl.NewSNSError(spl.InvalidAAAARecord, "the record content must be a valid IP address", nil)
		}

		if len(ip) != 16 {
			return nil, spl.NewSNSError(spl.InvalidAAAARecord, "The record content must be 16 bytes long", nil)
		}

		return ip, nil

	} else if record == types.Background {
		out, err := solana.PublicKeyFromBase58(str)
		if err != nil {
			return nil, spl.NewSNSError("", "the record content must be a valid PublicKey", err)
		}
		return out.Bytes(), nil
	}

	return nil, spl.NewSNSError(spl.InvalidRecordInput, "The provided record data is invalid", nil)
}

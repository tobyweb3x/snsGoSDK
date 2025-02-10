package recordv2

import (
	"encoding/hex"
	spl "snsGoSDK/spl"
	"snsGoSDK/types"
	"strings"

	"net"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/gagliardetto/solana-go"
	"golang.org/x/net/idna"
)

func SerializeRecordv2Contents(content string, record types.Record) ([]byte, error) {
	var err error
	if _, ok := UTF8Encoded[record]; ok {
		if record == types.CNAME || record == types.TXT {
			if content, err = idna.ToASCII(content); err != nil {
				return nil, nil
			}
			content = strings.TrimPrefix(content, "xn--")
		}
		return []byte(content), nil

	} else if record == types.SOL {
		out, err := solana.PublicKeyFromBase58(content)
		if err != nil {
			return nil, err
		}
		return out.Bytes(), nil

	} else if _, ok := EVMRecords[record]; ok {
		if !strings.HasPrefix(content, "0x") {
			return nil, spl.NewSNSError(spl.InvalidEvmAddress, "the record content must start with `0x`", nil)
		}
		content = strings.TrimPrefix(content, "0x")

		return hex.DecodeString(content)

	} else if record == types.Injective {
		hrp, decoded, err := bech32.Decode(content)
		if err != nil {
			return nil, spl.NewSNSError(spl.InvalidInjectiveAddress, "the record content must be a valid bech32 string", err)
		}
		if hrp != "inj" {
			return nil, spl.NewSNSError(spl.InvalidInjectiveAddress, "the record content must start with `inj`", nil)
		}
		if len(decoded) != 20 {
			return nil, spl.NewSNSError(spl.InvalidInjectiveAddress, "the record content must be 20 bytes long", nil)
		}

		return decoded, nil

	} else if record == types.A {
		ip := net.ParseIP(content)
		if ip == nil {
			return nil, spl.NewSNSError(spl.InvalidARecord, "the record content must be a valid IP address", nil)
		}
		if len(ip) != 4 {
			return nil, spl.NewSNSError(spl.InvalidARecord, "The record content must be 4 bytes long", nil)
		}
		return ip, nil

	} else if record == types.AAAA {
		ip := net.ParseIP(content)
		if ip == nil {
			return nil, spl.NewSNSError(spl.InvalidAAAARecord, "the record content must be a valid IP address", nil)
		}
		if len(ip) != 16 {
			return nil, spl.NewSNSError(spl.InvalidAAAARecord, "The record content must be 16 bytes long", nil)
		}
		return ip, nil
	}

	return nil, spl.NewSNSError(spl.InvalidRecordInput, "The record content is malformed", nil)
}

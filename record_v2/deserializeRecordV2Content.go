package recordv2

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/gagliardetto/solana-go"
	"golang.org/x/net/idna"
)

func DeserializeRecordV2Content(content []byte, record types.Record) (string, error) {
	var err error
	_, ok := UTF8Encoded[record]
	if ok {
		decoded := string(content)
		if record == types.CNAME || record == types.TXT {
			if decoded, err = idna.ToASCII(decoded); err != nil {
				return "", err
			}
			return strings.TrimPrefix(decoded, "xn--"), nil
		}
		return decoded, nil
	} else if record == types.SOL {
		if len(content) < 32 {
			return "", errors.New("content of byte slice is of length less than 32")
		}
		return solana.PublicKeyFromBytes(content).String(), nil
	} else if _, ok = EVMRecords[record]; ok {
		return fmt.Sprintf("0x%s", hex.EncodeToString(content)), nil
	} else if record == types.Injective {
		if content, err = bech32.ConvertBits(content, 8, 5, false); err != nil {
			return "", fmt.Errorf("err converting to bech32 base5: error: %v", err)
		}
		return bech32.Encode("inj", content)

	} else if record == types.A || record == types.AAAA {
		var ip net.IP
		if record == types.A {
			if ip = net.IP(content); ip.To4() == nil {
				return "", errors.New("invalid ip addr")
			}
			return ip.To4().String(), nil
		}

		if ip = net.IP(content); ip.To16() == nil {
			return "", errors.New("invalid ip addr")
		}

		return ip.To16().String(), nil
	}

	return "", spl.NewSNSError(spl.InvalidRecordData, "the record content is malformed", nil)
}

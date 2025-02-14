package resolve

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"snsGoSDK/nft"
	"snsGoSDK/record"
	recordV2 "snsGoSDK/record_v2"
	snsRecord "snsGoSDK/sns-record"
	"snsGoSDK/spl"
	"snsGoSDK/types"
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type AllowPDA string

const (
	AllowPDAAny   AllowPDA = "any"
	AllowPDATrue  AllowPDA = "true"
	AllowPDAFalse AllowPDA = "false"
)

type ResolveConfig struct {
	AllowPda   AllowPDA
	ProgramIDs []solana.PublicKey
}

func Resolve(
	conn *rpc.Client,
	domain string,
	config ResolveConfig,
) (solana.PublicKey, error) {

	out, err := utils.GetDomainKeySync(domain, types.VersionUnspecified)
	if err != nil {
		return solana.PublicKey{}, err
	}
	var nftRecord nft.NftRecord
	nftRecordKey, _, err := nftRecord.FindKey(out.PubKey,
		spl.NameTokenizerID)
	if err != nil {
		return solana.PublicKey{}, err
	}
	solRecordv1Key, err := record.GetRecordKey(domain, types.SOL)
	if err != nil {
		return solana.PublicKey{}, err
	}
	solRecord2Key, err := recordV2.GetRecordV2Key(domain, types.SOL)
	if err != nil {
		return solana.PublicKey{}, err
	}

	out2, err := conn.GetMultipleAccounts(
		context.TODO(),
		nftRecordKey,
		solRecordv1Key,
		solRecord2Key,
		out.PubKey,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}

	if len(out2.Value) <= 4 {
		return solana.PublicKey{}, errors.New("result list not complete")
	}

	nftRecordInfo, solRecordV1Info,
		solRecordV2Info, registryInfo := out2.Value[0], out2.Value[1],
		out2.Value[2], out2.Value[3]

	if registryInfo.Data == nil {
		return solana.PublicKey{},
			spl.NewSNSError(spl.DomainDoesNotExist,
				fmt.Sprintf("domain %s does not exist", domain), nil)
	}

	var registry spl.NameRegistryState
	if err := registry.Deserialize(registryInfo.Data.GetBinary()); err != nil {
		return solana.PublicKey{}, err
	}

	// If NFT record active -> NFT owner is the owner
	if nftRecordInfo.Data != nil {
		var nftRecord nft.NftRecord
		if err := borsh.Deserialize(&nftRecord, nftRecordInfo.Data.GetBinary()); err != nil {
			return solana.PublicKey{}, err
		}
		if nftRecord.Tag == nft.ActiveRecord {
			nftOwner, err := nft.RetrieveNftOwnerV2(conn, out.PubKey)
			if err != nil || nftOwner.IsZero() {
				return solana.PublicKey{},
					spl.NewSNSError(spl.CouldNotFindNftOwner, "", err)
			}
			return nftOwner, nil
		}
	}

	// Check SOL record V2
	if solRecordV2Info.Data != nil {
		var recordV2 snsRecord.Record
		if err := recordV2.Deserialize(solRecordV2Info.Data.GetBinary()); err != nil {
			return solana.PublicKey{}, err
		}
		stalenessId, err := recordV2.GetStalenessId()
		if err != nil {
			return solana.PublicKey{}, err
		}
		roaId, err := recordV2.GetRoAId()
		if err != nil {
			return solana.PublicKey{}, err
		}
		content, err := recordV2.GetContent()
		if err != nil {
			return solana.PublicKey{}, err
		}

		if len(content) != 32 {
			return solana.PublicKey{}, spl.NewSNSError(spl.RecordMalformed, "record is malformed", nil)
		}

		if recordV2.Header.RightOfAssociationValidation != uint16(snsRecord.Solana) ||
			recordV2.Header.StalenessValidation != uint16(snsRecord.Solana) {
			return solana.PublicKey{},
				spl.NewSNSError(spl.WrongValidation, "record has wrong validation", nil)
		}

		skipFlag := false
		if r := slices.Compare(stalenessId, registry.Owner.Bytes()); r != 0 {
			skipFlag = true
		}

		if !skipFlag {
			if r := slices.Compare(roaId, content); r == 0 {
				return solana.PublicKeyFromBytes(content), nil
			}

			return solana.PublicKey{},
				spl.NewSNSError(spl.InvalidRoA,
					fmt.Sprintf("the RoA Id should be %s but is %s",
						solana.PublicKeyFromBytes(content).String(), solana.PublicKeyFromBytes(roaId).String()), nil)
		}
	}

	// Check SOL record V1
	if solRecordV1Info.Data != nil &&
		len(solRecordV1Info.Data.GetBinary()) > spl.HEADER_LEN+32 {
		var expectedBuffer bytes.Buffer
		expectedBuffer.Write(solRecordV1Info.Data.GetBinary()[spl.HEADER_LEN : spl.HEADER_LEN+32])
		expectedBuffer.Write(solRecordv1Key.Bytes())

		expectedHex := hex.EncodeToString(expectedBuffer.Bytes())
		expected := []byte(expectedHex)

		valid := record.CheckSolRecord(
			expected,
			solRecordV1Info.Data.GetBinary()[spl.HEADER_LEN+32:spl.HEADER_LEN+solana.SignatureLength],
			registry.Owner,
		)
		if valid {
			return solana.PublicKeyFromBytes(
				solRecordV1Info.Data.GetBinary()[spl.HEADER_LEN : spl.HEADER_LEN+32],
			), nil
		}
	}

	// Check if the registry owner is a PDA
	if !solana.IsOnCurve(registry.Owner.Bytes()) {
		if config.AllowPda == AllowPDAAny {
			return registry.Owner, nil
		}

		if config.AllowPda == AllowPDATrue {
			ownerInfo, err := conn.GetAccountInfo(
				context.TODO(),
				registry.Owner,
			)
			if err != nil {
				return solana.PublicKey{}, err
			}

			isAllowed := false
			for _, programID := range config.ProgramIDs {
				if ownerInfo.Value.Owner.Equals(programID) {
					isAllowed = true
					break
				}
			}

			if isAllowed {
				return registry.Owner, nil
			}

			return solana.PublicKey{}, spl.NewSNSError(
				spl.PdaOwnerNotAllowed,
				fmt.Sprintf("the Program %s is not allowed", ownerInfo.Value.Owner.String()),
				nil,
			)
		}
	}

	return registry.Owner, nil
}

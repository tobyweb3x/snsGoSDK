package snsRecord

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

const NameRegistryLen = 96

type Validation int

const (
	None Validation = iota
	Solana
	Ethereum
	UnverifiedSolana
)

func getValidationLength(v Validation) (int, error) {
	switch v {
	case None:
		return 0, nil
	case Solana:
		return 32, nil
	case Ethereum:
		return 20, nil
	case UnverifiedSolana:
		return 32, nil
	default:
		return 0, fmt.Errorf("invalid validation type")
	}
}

const RecordHeaderLen = 8

type RecordHeader struct {
	StalenessValidation          uint16
	RightOfAssociationValidation uint16
	ContentLength                uint32
}

func NewRecordHeader(
	stalenessValidation uint16,
	rightOfAssociationValidation uint16,
	contentLength uint32,
) *RecordHeader {
	return &RecordHeader{
		StalenessValidation:          stalenessValidation,
		RightOfAssociationValidation: rightOfAssociationValidation,
		ContentLength:                contentLength,
	}
}

func (rh *RecordHeader) Retrieve(conn *rpc.Client, key solana.PublicKey) error {
	out, err := conn.GetAccountInfo(context.TODO(), key)
	if err != nil || out == nil || out.Value == nil {
		return fmt.Errorf("Record header account not found: %w", err)
	}

	return borsh.Deserialize(rh, out.Value.Data.GetBinary()[NameRegistryLen:NameRegistryLen+RecordHeaderLen])
}

type Record struct {
	Header RecordHeader
	Data   []byte
}

func NewRecord(
	header RecordHeader,
	data []byte,
) *Record {
	return &Record{
		Header: header,
		Data:   data,
	}
}

func (r *Record) Deserialize(buff []byte) error {
	header, offset := RecordHeader{}, NameRegistryLen
	err := borsh.Deserialize(&header, buff[offset:offset+RecordHeaderLen])
	if err != nil {
		return err
	}
	data := buff[offset+RecordHeaderLen:]
	r.Header = header
	r.Data = data
	return nil
}

func (r *Record) Retrieve(conn *rpc.Client, key solana.PublicKey) error {
	info, err := conn.GetAccountInfo(context.TODO(), key)
	if err != nil {
		return err
	}
	if info == nil || info.Value == nil || info.Value.Data == nil || info.Value.Data.GetBinary() == nil {
		return errors.New("record header account not found")
	}

	return r.Deserialize(info.Value.Data.GetBinary())
}
func (r Record) RetrieveBatch(conn *rpc.Client, key []solana.PublicKey) ([]Record, error) {
	// batchSize := 100
	info, err := conn.GetMultipleAccounts(context.TODO(), key...)
	if err != nil {
		return nil, err
	}

	if info == nil || info.Value == nil {
		return nil, errors.New("record header for accounts not found")
	}

	accs := info.Value
	if len(key) != len(accs) {
		return nil, errors.New("incomplete account data was returned")
	}

	records := make([]Record, len(accs))
	for i, v := range accs {
		if v.Data == nil {
			records[i] = Record{}
			continue
		}
		var r Record
		if err := r.Deserialize(v.Data.GetBinary()); err != nil {
			records[i] = Record{}
			continue
		}
		records[i] = r
	}

	return records, nil
}

func (r *Record) GetContent() ([]byte, error) {
	a, err := getValidationLength(Validation(r.Header.StalenessValidation))
	if err != nil {
		return nil, err
	}
	b, err := getValidationLength(Validation(r.Header.RightOfAssociationValidation))
	if err != nil {
		return nil, err
	}
	startOffset := a + b

	if len(r.Data) < startOffset {
		return nil, fmt.Errorf("data in Record is too short: expected at least %d bytes, got %d", startOffset, len(r.Data))
	}

	return r.Data[startOffset:], nil
}

func (r *Record) GetStalenessId() ([]byte, error) {
	endOffset, err := getValidationLength(Validation(r.Header.StalenessValidation))
	if err != nil {
		return nil, err
	}
	if len(r.Data) < endOffset {
		return nil, fmt.Errorf("invalid staleness validation length")
	}
	return r.Data[:endOffset], nil
}

func (r *Record) GetRoAId() ([]byte, error) {
	startOffset, err := getValidationLength(Validation(r.Header.StalenessValidation))
	if err != nil {
		return nil, err
	}
	endOffset, err := getValidationLength(Validation(r.Header.RightOfAssociationValidation))
	if err != nil {
		return nil, err
	}

	finalOffset := startOffset + endOffset
	if len(r.Data) < finalOffset {
		return nil, fmt.Errorf("data in Record is too short: expected at least %d bytes, got %d", finalOffset, len(r.Data))
	}

	return r.Data[startOffset:finalOffset], nil
}

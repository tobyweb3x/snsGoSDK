package snsRecord

import (
	"context"
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
		return 0, fmt.Errorf("Invalid validation type")
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
	offset := NameRegistryLen
	var header RecordHeader
	err := borsh.Deserialize(&header, buff[offset:offset+RecordHeaderLen])
	if err != nil {
		return err
	}
	data := buff[offset+RecordHeaderLen:]
	r.Header = header
	r.Data = data
	return nil
}

// func (r *Record) Retrieve(conn *rpc.Client, key solana.PublicKey) error {}

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
	return r.Data[startOffset:], nil

}

func (r *Record) GetStalenessId() ([]byte, error) {
	endOffset, err := getValidationLength(Validation(r.Header.StalenessValidation))
	if err != nil {
		return nil, err
	}
	if len(r.Data) < endOffset {
		return nil, fmt.Errorf("Invalid staleness validation length")
	}
	return r.Data[:endOffset], nil
}

func (r *Record) GetRoAId() ([]byte, error) {
	startOffset, err := getValidationLength(Validation(r.Header.StalenessValidation))
	if err != nil {
		return nil, err
	}
	endOfffset, err := getValidationLength(Validation(r.Header.RightOfAssociationValidation))
	if err != nil {
		return nil, err
	}

	return r.Data[startOffset : startOffset+endOfffset], nil
}

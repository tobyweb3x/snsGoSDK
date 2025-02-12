package spl

import "fmt"

// SNSError is a custom error type for Solana Name Service errors.
type SNSError struct {
	type_   ErrorType
	message string
}

func (e SNSError) Error() string {
	return fmt.Sprintf("[%s] %s", e.type_, e.message)
}

func (e SNSError) ErrorType() error {
	return e.type_
}

func (e SNSError) Message() string {
	return e.message
}

func NewSNSError(errorType ErrorType, message string, err error) SNSError {
	m := fmt.Errorf("%s: error : %s", message, err.Error())
	return SNSError{
		type_:   errorType,
		message: m.Error(),
	}
}

type TokenError error

var (
	//ErrTokenAccountNotFound: Thrown if an account is not found at the expected address, this error can be ignored.
	ErrTokenAccountNotFound TokenError = fmt.Errorf("ErrTokenAccountNotFound")

	//ErrTokenInvalidAccountOwner: Thrown if a program state account is not owned by the expected token program, this error can be ignored.
	ErrTokenInvalidAccountOwner TokenError = fmt.Errorf("ErrTokenInvalidAccountOwner")

	//ErrTokenInvalidAccountSize: Thrown if the byte length of an program state account doesn't match the expected size, this error can be ignored.
	ErrTokenInvalidAccountSize TokenError = fmt.Errorf("ErrTokenInvalidAccountSize")
)

var (
	// MintAccount has zero supply.
	ErrZeroMintSupply = fmt.Errorf("MintAccount has zero supply")
	// This error is intentionally ignored and program execution is continued.
	ErrIgnored = fmt.Errorf("this error is ignored")
)

type ErrorType string

func (e ErrorType) Error() string {
	return string(e)
}

const (
	SymbolNotFound                   ErrorType = "SymbolNotFound"
	InvalidSubdomain                 ErrorType = "InvalidSubdomain"
	FavouriteDomainNotFound          ErrorType = "FavouriteDomainNotFound"
	MissingParentOwner               ErrorType = "MissingParentOwner"
	U32Overflow                      ErrorType = "U32Overflow"
	InvalidBufferLength              ErrorType = "InvalidBufferLength"
	U64Overflow                      ErrorType = "U64Overflow"
	NoRecordData                     ErrorType = "NoRecordData"
	InvalidRecordData                ErrorType = "InvalidRecordData"
	UnsupportedRecord                ErrorType = "UnsupportedRecord"
	InvalidEvmAddress                ErrorType = "InvalidEvmAddress"
	InvalidInjectiveAddress          ErrorType = "InvalidInjectiveAddress"
	InvalidARecord                   ErrorType = "InvalidARecord"
	InvalidAAAARecord                ErrorType = "InvalidAAAARecord"
	InvalidRecordInput               ErrorType = "InvalidRecordInput"
	InvalidSignature                 ErrorType = "InvalidSignature"
	AccountDoesNotExist              ErrorType = "AccountDoesNotExist"
	MultipleRegistries               ErrorType = "MultipleRegistries"
	InvalidReverseTwitter            ErrorType = "InvalidReverseTwitter"
	NoAccountData                    ErrorType = "NoAccountData"
	InvalidInput                     ErrorType = "InvalidInput"
	InvalidDomain                    ErrorType = "InvalidDomain"
	InvalidCustomBg                  ErrorType = "InvalidCustomBackground"
	UnsupportedSignature             ErrorType = "UnsupportedSignature"
	RecordDoestNotSupportGuardianSig ErrorType = "RecordDoestNotSupportGuardianSig"
	RecordIsNotSigned                ErrorType = "RecordIsNotSigned"
	UnsupportedSignatureType         ErrorType = "UnsupportedSignatureType"
	InvalidSolRecordV2               ErrorType = "InvalidSolRecordV2"
	MissingVerifier                  ErrorType = "MissingVerifier"
	PythFeedNotFound                 ErrorType = "PythFeedNotFound"
	InvalidRoA                       ErrorType = "InvalidRoA"
	InvalidPda                       ErrorType = "InvalidPda"
	InvalidParrent                   ErrorType = "InvalidParrent"
	NftRecordNotFound                ErrorType = "NftRecordNotFound"
	PdaOwnerNotAllowed               ErrorType = "PdaOwnerNotAllowed"
	DomainDoesNotExist               ErrorType = "DomainDoesNotExist"
	RecordMalformed                  ErrorType = "RecordMalformed"
	CouldNotFindNftOwner             ErrorType = "CouldNotFindNftOwner"
	WrongValidation                  ErrorType = "WrongValidation"
)

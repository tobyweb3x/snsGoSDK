package spl_name_services

import "fmt"

// SNSError is a custom error type for Solana Name Service errors.
type SNSError error

var (
	// SNSError: InvalidInput
	ErrInvalidInput SNSError = fmt.Errorf("SNSError: InvalidInput")
	// SNSError: AccountDoesNotExist
	ErrAccountDoesNotExist SNSError = fmt.Errorf("SNSError: AccountDoesNotExist")
	// SNSError: ErrNoAccountData
	ErrNoAccountData SNSError = fmt.Errorf("SNSError: ErrNoAccountData")
)

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

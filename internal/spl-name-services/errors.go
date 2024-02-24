package spl_name_services

import "fmt"

type SNSError error

var (
	ErrInvalidInput        SNSError = fmt.Errorf("SNSError: InvalidInput")
	ErrAccountDoesNotExist SNSError = fmt.Errorf("SNSError: AccountDoesNotExist")
	ErrNoAccountData       SNSError = fmt.Errorf("SNSError: ErrNoAccountData")
)

type TokenError error

var (
	//ErrTokenAccountNotFound: Thrown if an account is not found at the expected address
	ErrTokenAccountNotFound TokenError = fmt.Errorf("ErrTokenAccountNotFound")

	//ErrTokenInvalidAccountOwner: Thrown if a program state account is not owned by the expected token program
	ErrTokenInvalidAccountOwner TokenError = fmt.Errorf("ErrTokenInvalidAccountOwner")

	ErrTokenInvalidAccountSize TokenError = fmt.Errorf("ErrTokenInvalidAccountSize")
)

var (
	ErrZeroMintSupply = fmt.Errorf("MintAccount has zero supply")
	ErrIgnored        = fmt.Errorf("this error is ignored")
)

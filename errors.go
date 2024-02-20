package main

import "fmt"

type SNSError error

var (
	ErrInvalidInput        SNSError = fmt.Errorf("SNSError: InvalidInput")
	ErrAccountDoesNotExist SNSError = fmt.Errorf("SNSError: AccountDoesNotExist")
)

type TokenError error

var (
	//TokenAccountNotFoundError: Thrown if an account is not found at the expected address
	TokenAccountNotFoundError TokenError = fmt.Errorf("TokenAccountNotFoundError")

	//TokenInvalidAccountOwnerError: Thrown if a program state account is not owned by the expected token program
	TokenInvalidAccountOwnerError TokenError = fmt.Errorf("TokenInvalidAccountOwnerError")

	TokenInvalidAccountSizeError TokenError = fmt.Errorf("TokenInvalidAccountSizeError")
	ErrZeroMintSupply                       = fmt.Errorf("MintAccount has zero supply")
)

package main

import "fmt"

type SNSError error 

var (
	ErrInvalidInput SNSError = fmt.Errorf("SNSError: InvalidInput")
)

func s() {
	ErrInvalidInput.Error() 
}

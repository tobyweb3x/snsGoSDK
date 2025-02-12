package snsRecord

type Validation int

const (
	None Validation = iota
	Solana
	Ethereum
	UnverifiedSolana
)

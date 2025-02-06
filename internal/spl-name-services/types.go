package spl_name_services

import (
	"github.com/gagliardetto/solana-go"
)

// DomainKeyResult is a helper struct for GetDomainKeySync.
type DomainKeyResult struct {
	PubKey      solana.PublicKey
	Parent      solana.PublicKey
	Hashed      []byte
	IsSub       bool
	IsSubRecord bool
}

// deriveResult is a helper struct for GetDomainKeySync.
type deriveResult struct {
	PubKey solana.PublicKey
	Hashed []byte
}

// RecordVersion type for spl-name-services.
type RecordVersion uint8

const (
	V1 RecordVersion = 1
	V2 RecordVersion = 2
)

const (
	HashPrefix = "SPL Name Service"
	HEADER_LEN = 96
)

var (
	/*
		The Solana Name Service program ID.
			NameProgramID = solana.MustPublicKeyFromBase58("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")
	*/
	NameProgramID = solana.MustPublicKeyFromBase58("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")

	/*
		Central State.
			CentralStateSNSRecords = solana.MustPublicKeyFromBase58("2pMnqHvei2N5oDcVGCRdZx48gqti199wr5CsyTTafsbo")

	*/
	CentralStateSNSRecords = solana.MustPublicKeyFromBase58("2pMnqHvei2N5oDcVGCRdZx48gqti199wr5CsyTTafsbo")

	/*
		The ".sol" TLD.
			RootDomainAccount = solana.MustPublicKeyFromBase58("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx")
	*/
	RootDomainAccount = solana.MustPublicKeyFromBase58("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx")

	/*
		NoPublicKeyArg is an alias for:
			solana.PublicKey{}
		so you compare with the equality operator.
	*/
	NoPublickKeyArg = solana.PublicKey{}

	/*
		Address of the SPL Token program.
			TOKEN_PROGRAM_ID = solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	*/
	TOKEN_PROGRAM_ID = solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

	/*
		The reverse look up class.
			REVERSE_LOOKUP_CLASS = solana.MustPublicKeyFromBase58("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")
	*/
	REVERSE_LOOKUP_CLASS = solana.MustPublicKeyFromBase58("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")
)

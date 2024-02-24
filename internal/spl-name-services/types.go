package spl_name_services

import (
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/rpc"
)

// DomainKeyResult is a helper struct for GetDomainKeySync.
type DomainKeyResult struct {
	PubKey      common.PublicKey
	Parent      common.PublicKey
	Hashed      []byte
	IsSub       bool
	IsSubRecord bool
}

// deriveResult is a helper struct for GetDomainKeySync.
type deriveResult struct {
	PubKey common.PublicKey
	Hashed []byte
}

// RecordVersion type for spl-name-services.
type RecordVersion int8

const (
	/*
		RecordVersion
			RecordVersion1 RecordVersion = 1
	*/
	RecordVersion1 RecordVersion = 1
	/*
		RecordVersion
			RecordVersion2 RecordVersion = 2
	*/
	RecordVersion2 RecordVersion = 2
)

const (
	HashPrefix = "SPL Name Service"
	HEADER_LEN = 96
)

const (
	/*
		NoCommitmentArg is an alias for:
			NoCommitmentArg rpc.Commitment = ""
	*/
	NoCommitmentArg rpc.Commitment = ""
)

var (
	/*
		The Solana Name Service program ID.
			NameProgramID = common.PublicKeyFromString("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")
	*/
	NameProgramID = common.PublicKeyFromString("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")

	/*
		Central State.
			CentralStateSNSRecords = common.PublicKeyFromString("2pMnqHvei2N5oDcVGCRdZx48gqt i199v5CsyTTafsbo")
	*/
	CentralStateSNSRecords = common.PublicKeyFromString("2pMnqHvei2N5oDcVGCRdZx48gqt i199v5CsyTTafsbo")

	/*
		The ".sol" TLD.
			RootDomainAccount = common.PublicKeyFromString("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx")
	*/
	RootDomainAccount = common.PublicKeyFromString("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx")

	/*
		NoPublicKeyArg is an alias for:
			common.PublicKey{}
		so you compare with the equality operator.
	*/
	NoPublickKeyArg = common.PublicKey{}

	/*
		MINT_PREFIX.
			MINT_PREFIX = []byte("tokenized_name")
	*/
	MINT_PREFIX = []byte("tokenized_name")

	/*
		NAME_TOKENIZER_ID.
			NAME_TOKENIZER_ID = common.PublicKeyFromString("nftD3vbNkNqfj2Sd3HZwbpw4BxxKWr4AjGb9X38JeZk")
	*/
	NAME_TOKENIZER_ID = common.PublicKeyFromString("nftD3vbNkNqfj2Sd3HZwbpw4BxxKWr4AjGb9X38JeZk")

	/*
		Address of the SPL Token program.
			TOKEN_PROGRAM_ID = common.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	*/
	TOKEN_PROGRAM_ID = common.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

	/*
		The reverse look up class.
			REVERSE_LOOKUP_CLASS = common.PublicKeyFromString("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")
	*/
	REVERSE_LOOKUP_CLASS = common.PublicKeyFromString("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")
)

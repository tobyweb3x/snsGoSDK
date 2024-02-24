package spl_name_services

import (
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/rpc"
)

type DomainKeyResult struct {
	PubKey      common.PublicKey
	Parent      common.PublicKey
	Hashed      []byte
	IsSub       bool
	IsSubRecord bool
}

type deriveResult struct {
	PubKey common.PublicKey
	Hashed []byte
}

type RecordVersion int8

const (
	RecordVersion1 RecordVersion = 1
	RecordVersion2 RecordVersion = 2
)

const (
	HashPrefix = "SPL Name Service"
	HEADER_LEN = 96
)

const (
	NoCommitmentArg rpc.Commitment = ""
)

var (
	NameProgramID          = common.PublicKeyFromString("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")
	CentralStateSNSRecords = common.PublicKeyFromString("2pMnqHvei2N5oDcVGCRdZx48gqt i199v5CsyTTafsbo")
	RootDomainAccount      = common.PublicKeyFromString("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx")
	/*
		NoPublicKeyArg is an alias for:
			common.PublicKey{}
		so you compare with the equality operator
	*/
	NoPublickKeyArg   = common.PublicKey{}
	MINT_PREFIX       = []byte("tokenized_name")
	NAME_TOKENIZER_ID = common.PublicKeyFromString("nftD3vbNkNqfj2Sd3HZwbpw4BxxKWr4AjGb9X38JeZk")

	TOKEN_PROGRAM_ID     = common.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	REVERSE_LOOKUP_CLASS = common.PublicKeyFromString("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")
)

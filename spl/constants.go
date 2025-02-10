package spl_name_services

import "github.com/gagliardetto/solana-go"

const (
	HashPrefix = "SPL Name Service"
	HEADER_LEN = 96
)

var (

	// MIntPrefix.
	//  MIntPrefix = []byte("tokenized_name")
	MIntPrefix = []byte("tokenized_name")

	// The program ID of the Name Service program.
	//  NameTokenizerID = solana.MustPublicKeyFromBase58("nftD3vbNkNqfj2Sd3HZwbpw4BxxKWr4AjGb9X38JeZk")
	NameTokenizerID = solana.MustPublicKeyFromBase58("nftD3vbNkNqfj2Sd3HZwbpw4BxxKWr4AjGb9X38JeZk")

	// The `.twitter` TLD
	//  TwitterRootParentRegistryKey = solana.PublicKeyFromString("4YcexoW3r78zz16J2aqmukBLRwGq6rAvWzJpkYAXqebv")
	TwitterRootParentRegistryKey = solana.MustPublicKeyFromBase58("4YcexoW3r78zz16J2aqmukBLRwGq6rAvWzJpkYAXqebv")

	// The ".twitter" TLD authority.
	//  TwittwrVerificationAuthority = solana.PublicKeyFromString("FvPH7PrVrLGKPfqaf3xJodFTjZriqrAXXLTVWEorTFBi")
	TwittwrVerificationAuthority = solana.MustPublicKeyFromBase58("FvPH7PrVrLGKPfqaf3xJodFTjZriqrAXXLTVWEorTFBi")

	// The Registry program ID
	//  ResgistryProgramID = solana.MustPublicKeyFromBase58("jCebN34bUfdeUYJT13J1yG16XWQpt5PDx6Mse9GUqhR")
	ResgistryProgramID = solana.MustPublicKeyFromBase58("jCebN34bUfdeUYJT13J1yG16XWQpt5PDx6Mse9GUqhR")

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
		Address of the SPL Token program.
			TOKEN_PROGRAM_ID = solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	*/
	TOKEN_PROGRAM_ID = solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

	/*
		The reverse look up class.
			ReverseLookupClass = solana.MustPublicKeyFromBase58("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")
	*/
	ReverseLookupClass = solana.MustPublicKeyFromBase58("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")

	CentralState = ReverseLookupClass
)

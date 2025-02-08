package spl_name_services

import "github.com/gagliardetto/solana-go"

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
)

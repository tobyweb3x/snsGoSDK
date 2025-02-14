package spl

import "github.com/gagliardetto/solana-go"

const (
	HashPrefix = "SPL Name Service"
	NameRegistryStateHeaderLen = 96
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

	// The Solana Name Service program ID.
	// 	NameProgramID = solana.MustPublicKeyFromBase58("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")
	NameProgramID = solana.MustPublicKeyFromBase58("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")

	// Central State.
	//  CentralStateSNSRecords = solana.MustPublicKeyFromBase58("2pMnqHvei2N5oDcVGCRdZx48gqti199wr5CsyTTafsbo")
	CentralStateSNSRecords = solana.MustPublicKeyFromBase58("2pMnqHvei2N5oDcVGCRdZx48gqti199wr5CsyTTafsbo")

	// The ".sol" TLD.
	// 	RootDomainAccount = solana.MustPublicKeyFromBase58("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx")
	RootDomainAccount = solana.MustPublicKeyFromBase58("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx")

	// The reverse look up class.
	// 	ReverseLookupClass = solana.MustPublicKeyFromBase58("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")
	ReverseLookupClass = solana.MustPublicKeyFromBase58("33m47vH6Eav6jr5Ry86XjhRft2jRBLDnDgPSHoquXi2Z")

	// The Wolves Collection Metadata.
	//  WolvesCollectionMetadata = solana.MustPublicKeyFromBase58("72aLKvXeV4aansAQtxKymeXDevT5ed6sCuz9iN62ugPT")
	WolvesCollectionMetadata = solana.MustPublicKeyFromBase58("72aLKvXeV4aansAQtxKymeXDevT5ed6sCuz9iN62ugPT")
	CentralState = ReverseLookupClass

	// The MPL Token Metadata.
	//  MetaplexID = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
	MetaplexID = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

	// The Name Offers program ID.
	//  NameOffersID = solana.MustPublicKeyFromBase58("85iDfUvr3HJyLM2zcq5BXSiDvUWfw6cSE1FfNBo8Ap29")
	NameOffersID = solana.MustPublicKeyFromBase58("85iDfUvr3HJyLM2zcq5BXSiDvUWfw6cSE1FfNBo8Ap29")

	// USDC mint.
	//  USDCMint = solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	USDCMint = solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")

	// The Register program ID.
	//  RegisterProgramID = solana.MustPublicKeyFromBase58("jCebN34bUfdeUYJT13J1yG16XWQpt5PDx6Mse9GUqhR")
	RegisterProgramID = solana.MustPublicKeyFromBase58("jCebN34bUfdeUYJT13J1yG16XWQpt5PDx6Mse9GUqhR")

	/*
		REFERRERS = [...]solana.PublicKey{
			solana.MustPublicKeyFromBase58("3ogYncmMM5CmytsGCqKHydmXmKUZ6sGWvizkzqwT7zb1"), // Test wallet
			solana.MustPublicKeyFromBase58("DM1jJCkZZEwY5tmWbgvKRxsDFzXCdbfrYCCH1CtwguEs"), // 4Everland
			solana.MustPublicKeyFromBase58("ADCp4QXFajHrhy4f43pD6GJFtQLkdBY2mjS9DfCk7tNW"), // Bandit network
			solana.MustPublicKeyFromBase58("2XTgjw8yi1E3Etgj4CUyRD7Zk49gynH2U9gA5N2MY4NP"), // Altoscan
			solana.MustPublicKeyFromBase58("5PwNeqQPiygQks9R17jUAodZQNuhvCqqkrxSaeNE8qTR"), // Solscan
			solana.MustPublicKeyFromBase58("8kJqxAbqbPLGLMgB6FhLcnw2SiUEavx2aEGM3WQGhtJF"), // Domain Labs
			solana.MustPublicKeyFromBase58("HemvJzwxvVpWBjPETpaseAH395WAxb2G73MeUfjVkK1u"), // Solflare
			solana.MustPublicKeyFromBase58("7hMiiUtkH4StMPJxyAtvzXTUjecTniQ8czkCPusf5eSW"), // Solnames
			solana.MustPublicKeyFromBase58("DGpjHo4yYA3NgHvhHTp3XfBFrESsx1DnhfTr8D881ZBM"), // Brave
			solana.MustPublicKeyFromBase58("7vWSqSw1eCXZXXUubuHWssXELNQ8MLaDgAs2ErEfCKxn"), // 585.eth
			solana.MustPublicKeyFromBase58("5F6gcdzpw7wUjNEugdsD4aLJdEQ4Wt8d6E85vaQXZQSJ"), // wdotsol
			solana.MustPublicKeyFromBase58("XEy9o73JBN2pEuN7aspe8mVLaWbL4ozjJs1tNRxx8bL"), // GoDID
			solana.MustPublicKeyFromBase58("D5cLoAGjNTHKU1UGv2bYwbnyRoGTMe3sbpLtJW3fRq91"), // SuiNS
			solana.MustPublicKeyFromBase58("FePcCmrr7vgjeFXcXtJHqShSXydaTrga2wfHRt9RrYvP"), // Nansen
		}
	*/
	REFERRERS = []solana.PublicKey{
		solana.MustPublicKeyFromBase58("3ogYncmMM5CmytsGCqKHydmXmKUZ6sGWvizkzqwT7zb1"), // Test wallet
		solana.MustPublicKeyFromBase58("DM1jJCkZZEwY5tmWbgvKRxsDFzXCdbfrYCCH1CtwguEs"), // 4Everland
		solana.MustPublicKeyFromBase58("ADCp4QXFajHrhy4f43pD6GJFtQLkdBY2mjS9DfCk7tNW"), // Bandit network
		solana.MustPublicKeyFromBase58("2XTgjw8yi1E3Etgj4CUyRD7Zk49gynH2U9gA5N2MY4NP"), // Altoscan
		solana.MustPublicKeyFromBase58("5PwNeqQPiygQks9R17jUAodZQNuhvCqqkrxSaeNE8qTR"), // Solscan
		solana.MustPublicKeyFromBase58("8kJqxAbqbPLGLMgB6FhLcnw2SiUEavx2aEGM3WQGhtJF"), // Domain Labs
		solana.MustPublicKeyFromBase58("HemvJzwxvVpWBjPETpaseAH395WAxb2G73MeUfjVkK1u"), // Solflare
		solana.MustPublicKeyFromBase58("7hMiiUtkH4StMPJxyAtvzXTUjecTniQ8czkCPusf5eSW"), // Solnames
		solana.MustPublicKeyFromBase58("DGpjHo4yYA3NgHvhHTp3XfBFrESsx1DnhfTr8D881ZBM"), // Brave
		solana.MustPublicKeyFromBase58("7vWSqSw1eCXZXXUubuHWssXELNQ8MLaDgAs2ErEfCKxn"), // 585.eth
		solana.MustPublicKeyFromBase58("5F6gcdzpw7wUjNEugdsD4aLJdEQ4Wt8d6E85vaQXZQSJ"), // wdotsol
		solana.MustPublicKeyFromBase58("XEy9o73JBN2pEuN7aspe8mVLaWbL4ozjJs1tNRxx8bL"),  // GoDID
		solana.MustPublicKeyFromBase58("D5cLoAGjNTHKU1UGv2bYwbnyRoGTMe3sbpLtJW3fRq91"), // SuiNS
		solana.MustPublicKeyFromBase58("FePcCmrr7vgjeFXcXtJHqShSXydaTrga2wfHRt9RrYvP"), // Nansen
	}

	// The vault owner.
	//  VaultOwner = solana.MustPublicKeyFromBase58("5D2zKog251d6KPCyFyLMt3KroWwXXPWSgTPyhV22K2gR")
	VaultOwner     = solana.MustPublicKeyFromBase58("5D2zKog251d6KPCyFyLMt3KroWwXXPWSgTPyhV22K2gR")
	PythMappingAcc = solana.MustPublicKeyFromBase58("AHtgzX45WTKfkPG53L6WYhGEXwQkN1BVknET3sVsLL8J")

	// The default Pyth push program ID.
	//  DefaultPythPushProgram = solana.MustPublicKeyFromBase58("pythWSnswVUd12oZpeFP8e9CVaEqJg25g1Vtc2biRsT")
	DefaultPythPushProgram = solana.MustPublicKeyFromBase58("pythWSnswVUd12oZpeFP8e9CVaEqJg25g1Vtc2biRsT")

	PYTHFeeds = map[string]struct {
		Price   string
		Product string
	}{
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": {
			Price:   "Gnt27xtC473ZT2Mw5u8wZ68Z3gULkSTb5DuxJy7eJotD",
			Product: "8GWTTbNiXdmyZREXbjsZBmCRuzdPrW55dnZGDkTRjWvb",
		},
		"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB": {
			Price:   "3vxLXJqLqF3JG5TCbYycbKWRBbCJQLxQmBGCkyqEEefL",
			Product: "Av6XyAMJnyi68FdsKSPYgzfXGjYrrt6jcAMwtvzLCqaM",
		},
		"So11111111111111111111111111111111111111112": {
			Price:   "H6ARHf6YXhGYeQfUzQNGk6rDNnLBQKrenN712K4AQJEG",
			Product: "ALP8SdU9oARYVLgLR7LrqMNCYBnhtnQz1cj6bwgwQmgj",
		},
		"EchesyfXePKdLtoiZSL8pBe8Myagyy8ZRqsACNCFGnvp": {
			Price:   "ETp9eKXVv1dWwHSpsXRUuXHmw24PwRkttCGVgpZEY9zF",
			Product: "HyEB4Goiv7QyfFStaBn49JqQzSTV1ybtVikwsMLH1f2M",
		},
		"mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So": {
			Price:   "E4v1BBgoso9s64TQvmyownAVJbhbEPGyzA3qn4n46qj9",
			Product: "BS2iAqT67j8hA9Jji4B8UpL3Nfw9kwPfU5s4qeaf1e7r",
		},
		"DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263": {
			Price:   "8ihFLu5FimgTQ1Unh4dVyEHUGodJ5gJQCrQf4KUVB9bN",
			Product: "FerFD54J6RgmQVCR5oNgpzXmz8BW2eBNhhirb1d5oifo",
		},
		"EPeUFDgHRxs9xxEPVaL6kfGQvCon7jmAWKVUHuux1Tpz": {
			Price:   "AbMTYZ82Xfv9PtTQ5e1fJXemXjzqEEFHP3oDLRTae6yz",
			Product: "8xTEctXKo6Xo3EzNhSNp4TUe8mgfwWFbDUXJhuubvrKx",
		},
		"HZ1JovNiVvGrGNiiYvEozEVgZ58xaU3RKwX8eACQBCt3": {
			Price:   "nrYkQQQur7z8rYTST3G9GqATviK5SxTDkrqd21MW6Ue",
			Product: "AiQB4WngNPKDe3iWAwZmMzbULDAAfUD6Sr1knfZNJj3y",
		},
		"bSo13r4TkiE4KumL71LsHTPpL2euBYLFx6h9HP3piy1": {
			Price:   "AFrYBhb5wKQtxRS9UA9YRS4V3dwFm7SqmS6DHKq6YVgo",
			Product: "3RtUHQR2LQ7su5R4zWwjupx72sWRGvLA4cFmnbHnT9M7",
		},
		"6McPRfPV6bY1e9hLxWyG54W9i9Epq75QBvXg2oetBVTB": {
			Price:   "9EdtbaivHQYA4Nh3XzGR6DwRaoorqXYnmpfsnFhvwuVj",
			Product: "5Q5kyCVzssrGMd2BniSdVeRwjNWrGGrFhMrgGt4zURyA",
		},
	}

	PYTHPullFeeds = map[string][]byte{
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": {
			234, 160, 32, 198, 28, 196, 121, 113, 40, 19, 70, 28, 225, 83, 137, 74,
			150, 166, 192, 11, 33, 237, 12, 252, 39, 152, 209, 249, 169, 233, 201, 74,
		},
		"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB": {
			43, 137, 185, 220, 143, 223, 159, 52, 112, 154, 91, 16, 107, 71, 47, 15,
			57, 187, 108, 169, 206, 4, 176, 253, 127, 46, 151, 22, 136, 226, 229, 59,
		},
		"So11111111111111111111111111111111111111112": {
			239, 13, 139, 111, 218, 44, 235, 164, 29, 161, 93, 64, 149, 209, 218, 57,
			42, 13, 47, 142, 208, 198, 199, 188, 15, 76, 250, 200, 194, 128, 181, 109,
		},
		"EchesyfXePKdLtoiZSL8pBe8Myagyy8ZRqsACNCFGnvp": {
			200, 6, 87, 183, 246, 243, 234, 194, 114, 24, 208, 157, 90, 78, 84, 228,
			123, 37, 118, 141, 159, 94, 16, 172, 21, 254, 44, 249, 0, 136, 20, 0,
		},
		"mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So": {
			194, 40, 154, 106, 67, 210, 206, 145, 198, 245, 92, 174, 195, 112, 244,
			172, 195, 138, 46, 212, 119, 245, 136, 19, 51, 76, 109, 3, 116, 159, 242,
			164,
		},
		"DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263": {
			114, 176, 33, 33, 124, 163, 254, 104, 146, 42, 25, 170, 249, 144, 16, 156,
			185, 216, 78, 154, 208, 4, 180, 210, 2, 90, 214, 245, 41, 49, 68, 25,
		},
		"EPeUFDgHRxs9xxEPVaL6kfGQvCon7jmAWKVUHuux1Tpz": {
			142, 134, 15, 183, 78, 96, 229, 115, 107, 69, 93, 130, 246, 11, 55, 40, 4,
			156, 52, 142, 148, 150, 26, 221, 95, 150, 27, 2, 253, 238, 37, 53,
		},
		"HZ1JovNiVvGrGNiiYvEozEVgZ58xaU3RKwX8eACQBCt3": {
			11, 191, 40, 233, 168, 65, 161, 204, 120, 143, 106, 54, 27, 23, 202, 7,
			45, 14, 163, 9, 138, 30, 93, 241, 195, 146, 45, 6, 113, 149, 121, 255,
		},
		"bSo13r4TkiE4KumL71LsHTPpL2euBYLFx6h9HP3piy1": {
			137, 135, 83, 121, 231, 15, 143, 186, 220, 23, 174, 243, 21, 173, 243,
			168, 213, 209, 96, 184, 17, 67, 85, 55, 224, 60, 151, 232, 170, 201, 125,
			156,
		},
		"6McPRfPV6bY1e9hLxWyG54W9i9Epq75QBvXg2oetBVTB": {
			122, 91, 193, 210, 181, 106, 208, 41, 4, 140, 214, 57, 100, 179, 173, 39,
			118, 234, 223, 129, 46, 220, 26, 67, 163, 20, 6, 203, 84, 191, 245, 146,
		},
	}
)

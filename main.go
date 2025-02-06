package main

import (
	"fmt"
	"log"
	spl "snsGoSDK/internal/spl-name-services"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {

	testStruct := struct {
		domainName []string
	}{
		[]string{
			// "bonfida.sol",
			"tobytobias.sol",
		},
	}

	for k, v := range testStruct.domainName {
		fmt.Printf("\n%d). domain - %s\n", k+1, v)

		res, err := spl.GetDomainKeySync(v, spl.V2)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("GetDomainKeySync --> %+v\n\n", res)
		}

		rpcClient := createConnection("https://mainnet.helius-rpc.com/?api-key=13af3657-7609-4ede-9305-6ea6c7a2243f")

		reverseLookUPStr, err := spl.ReverseLookup(rpcClient, res.PubKey, solana.PublicKey{})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("reverse lookup for PublicKey %s is ----> %s \n\n", res.PubKey.String(), reverseLookUPStr)
		}

		var nm spl.NameRegistryState
		res2, err := nm.Retrieve(rpcClient, res.PubKey)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("nm.Retrieve --> %+v, %+v\n", *res2.Registry, res2.NftOwner)
			// if v == "bonfida.sol" { // if you loop here, it sorta complains about two many requests, rate limiting
			// 	testGetAllDomains(rpcClient, res2.Registry.Owner)
			// 	testGetDomainKeysWithReverses(rpcClient, res2.Registry.Owner)
		}

		// a, err := spl.GetAllDomains(rpcClient, res2.Registry.Owner)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(a, len(a))

		// b, err := spl.GetAllRegisteredDomain(rpcClient)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(b, len(b))

	}

	fmt.Printf("\n\n")
	// for _, v := range []string{"@oluwatobialone"} {
	// 	testTwitterResolving(rpcClient, v)
	// }

}

func createConnection(connType string) *rpc.Client {

	return rpc.New(connType)
}

// func testTwitterResolving(rpcClient *rpc.Client, twitterHandle string) {
// 	a, err := twitter.GetTwitterRegistry(rpcClient, twitterHandle)
// 	if err != nil {
// 		fmt.Printf("Error: twitter.GetTwitterRegistry for handle %s: %v", twitterHandle, err)
// 		return
// 	}

// 	fmt.Printf("Public Wallet address associated to the Twitter handle %s ----> %s\n", twitterHandle, a.Registry.Owner.String())

// 	b, c, err := twitter.GetHandleAndRegistryKey(rpcClient, a.Registry.Owner)
// 	if err != nil {
// 		fmt.Printf("Error: twitter.GetHandleAndRegistryKey for handle %s: %v", twitterHandle, err)
// 		return
// 	}

// 	fmt.Printf("Domain key associated to the Twitter handle '%s' is %s\n", c, b.String())
// 	fmt.Printf("\n\n")
// }

func testGetAllDomains(conn *rpc.Client, pubkey solana.PublicKey) {
	fmt.Println("*****GetAllDomains*****")

	a, err := spl.GetAllDomains(conn, pubkey)
	if err != nil {
		log.Fatalln("err from GetAllDomain---", err)
		return
	}
	for _, c := range a {
		fmt.Println("---->", c)
	}
}

func testGetDomainKeysWithReverses(conn *rpc.Client, pubkey solana.PublicKey) {
	fmt.Println("*****GetDomainKeysWithReverses*****")

	a, err := spl.GetDomainKeysWithReverses(conn, pubkey)
	if err != nil {
		log.Fatalln("err from GetDomainKeysWithReverses---", err)
		return
	}
	for _, c := range a {
		fmt.Println(c.PubKey, "---->", c.Domain)
	}
}

package main

import (
	"fmt"
	"log"

	spl "snsGoSDK/internal/spl-name-services"
	"snsGoSDK/internal/twitter"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/rpc"
)

func main() {
	var (
		res              spl.DomainKeyResult
		res2             spl.RetrieveResult
		err              error
		reverseLookUPStr string
		nm               spl.NameRegistryState
	)

	rpcClient := createConnection(rpc.MainnetRPCEndpoint)

	testStruct := struct {
		domainName []string
	}{
		[]string{"bonfida.sol", "solana.sol", "01.sol", "dex.solana.sol", "dex.bonfida.sol", "wallet-guide-5.sol", "sub-0.wallet-guide-3.sol"},
	}

	for k, v := range testStruct.domainName {
		fmt.Printf("\n\n%d). domain - %s\n", k+1, v)

		fmt.Println("*****GetDomainKeySync*****")
		if res, err = spl.GetDomainKeySync(v, spl.RecordVersion2); err != nil {
			fmt.Println(err)
		} else {
			// fmt.Printf("pubkey for `%s` is -- %s AND THE STRUCT VALUES ARE %+v\n", v, res.PubKey.ToBase58(), res)
			fmt.Printf("pubkey for `%s` is ----> %s\n", v, res.PubKey.ToBase58())
		}

		fmt.Println("*****ReverseLookUp*****")
		if reverseLookUPStr, err = spl.ReverseLookup(rpcClient, res.PubKey); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("reverse lookup for PublicKey %s is ----> `%s` \n", res.PubKey.ToBase58(), reverseLookUPStr)
		}

		fmt.Println("*****NameStateRegistry.Retrieve*****")
		if res2, err = nm.Retrieve(rpcClient, res.PubKey); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Owner's PublicKey for domain `%s` ----> %s\n\n", v, res2.Registry.Owner)
			if v == "bonfida.sol" { // if you loop here, it sorta complains about two many requests, rate limiting
				testGetAllDomains(rpcClient, res2.Registry.Owner)
				testGetDomainKeysWithReverses(rpcClient, res2.Registry.Owner)
			}
		}

	}

	fmt.Printf("\n\n")
	for _, v := range []string{"@oluwatobialone"} {
		testTwitterResolving(rpcClient, v)
	}

}

func createConnection(connType string) *client.Client {
	return client.NewClient(connType)
}

func testTwitterResolving(rpcClient *client.Client, twitterHandle string) {
	a, err := twitter.GetTwitterRegistry(rpcClient, twitterHandle)
	if err != nil {
		fmt.Printf("Error: twitter.GetTwitterRegistry for handle %s: %v", twitterHandle, err)
		return
	}

	fmt.Printf("Public Wallet address associated to the Twitter handle %s ----> %s\n", twitterHandle, a.Registry.Owner.ToBase58())

	b, c, err := twitter.GetHandleAndRegistryKey(rpcClient, a.Registry.Owner)
	if err != nil {
		fmt.Printf("Error: twitter.GetHandleAndRegistryKey for handle %s: %v", twitterHandle, err)
		return
	}

	fmt.Printf("Domain key associated to the Twitter handle '%s' is %s\n", c, b.ToBase58())
	fmt.Printf("\n\n")
}

func testGetAllDomains(conn *client.Client, pubkey common.PublicKey) {
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

func testGetDomainKeysWithReverses(conn *client.Client, pubkey common.PublicKey) {
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

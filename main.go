package main

import (
	"fmt"
	"log"

	"snsGoSDK/internal/twitter"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
)

func main() {
	// 	var (
	// 		res DomainKeyResult
	// 		// publicKey        common.PublicKey
	// 		res2             RetrieveResult
	// 		err              error
	// 		reverseLookUPStr string
	// 		nm               NameRegistryState
	// 	)

	rpcClient := createConnection(rpc.MainnetRPCEndpoint)

	// 	/******************************************
	// 	 * Bonfida
	// 	 ******************************************/

	// 	testStruct := struct {
	// 		domainName []string
	// 	}{
	// 		[]string{"bonfida", "wallet-guide-5.sol", "sub-0.wallet-guide-3.sol"},
	// 	}

	// 	for k, v := range testStruct.domainName {

	// 		fmt.Printf(
	// 			`******************************************
	// * test %d, for domain --- %s
	// ******************************************`+"\n", k, v)

	// 		fmt.Println("*****GetDomainKeySync*****")

	// 		if res, err = GetDomainKeySync(v, RecordVersion2); err != nil {
	// 			fmt.Println(err)
	// 		} else {
	// 			fmt.Printf("pubkey for `%s` is -- %s AND THE STRUCT VALUES ARE %+v\n", v, res.PubKey.ToBase58(), res)
	// 		}
	// 		fmt.Println("*******************************************")

	// 		fmt.Println("*****ReverseLookUp*****")

	// 		if reverseLookUPStr, err = ReverseLookup(rpcClient, res.PubKey); err != nil {
	// 			fmt.Println(err)
	// 		} else {
	// 			fmt.Printf("reverse lookup for `%s` pubkey(%s) is = `%s` \n", v, res.PubKey.ToBase58(), reverseLookUPStr)
	// 		}
	// 		fmt.Println("*******************************************")

	// 		fmt.Println("*****NameStateRegistry.Retrieve*****")

	// 		if res2, err = nm.Retrieve(rpcClient, res.PubKey); err != nil {
	// 			fmt.Println(err)
	// 		} else {
	// 			fmt.Printf("owner for `%s` is %s\n", v, res2.Registry.Owner)
	// 		}
	// 		fmt.Println("*******************************************")
	// 		fmt.Printf("\n\n")

	// }

	a, err := twitter.GetTwitterRegistry(rpcClient, "@oluwatobialone")
	if err != nil {
		log.Fatal("err:", err)
	}
	fmt.Printf("public key associated to the Twitter handle @oluwatobialone is %s\n", a.Registry.Owner.ToBase58())

	b, c, err := twitter.GetHandleAndRegistryKey(rpcClient, a.Registry.Owner)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b, c)
}

func createConnection(connType string) *client.Client {
	return client.NewClient(connType)
}

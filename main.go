package main

import (
	"fmt"

	"github.com/blocto/solana-go-sdk/rpc"
)

func main() {
	var (
		res DomainKeyResult
		// publicKey        common.PublicKey
		res2             RetrieveResult
		err              error
		reverseLookUPStr string
		nm               NameRegistryState
	)

	rpcClient := createConnection(rpc.MainnetRPCEndpoint)

	/******************************************
	 * Bonfida
	 ******************************************/

	if res, err = GetDomainKeySync("bonfida", RecordVersion2); err != nil {
		fmt.Println("error was ", err)
	} else {
		fmt.Printf("pubkey for bonfida and struct values are %+v\n", res.PubKey.ToBase58())
	}

	if reverseLookUPStr, err = ReverseLookup(rpcClient, res.PubKey); err != nil {
		fmt.Println("reverseLookup: err was ", err)
	} else {
		fmt.Printf("reverse lookup for bonfida pubkey(%s) is = `%s` ", res.PubKey.ToBase58(), reverseLookUPStr)
	}

	if res2, err = nm.Retrieve(rpcClient, res.PubKey); err != nil {
		fmt.Println("Retrieve: error was ", err)
	} else {
		fmt.Println("owner for bonfida ", res2.NftOwner.ToBase58())
	}
	fmt.Printf("\n\n")

	/******************************************
	* dex.bonfida
	******************************************/

	if res, err = GetDomainKeySync("dex.bonfida", RecordVersion2); err != nil {
		fmt.Println("error was ", err)
	} else {
		fmt.Printf("pubkey for dex.bonfand struct values are %+v\nida ", res.PubKey.ToBase58())
	}

	if reverseLookUPStr, err = ReverseLookup(rpcClient, res.PubKey); err != nil {
		fmt.Println("reverseLookup: err was ", err)
	} else {
		fmt.Printf("reverse lookup for dex.bonfida pubkey(%s) is = `%s` ", res.PubKey.ToBase58(), reverseLookUPStr)
	}

	if res2, err = nm.Retrieve(rpcClient, res.PubKey); err != nil {
		fmt.Println("Retrieve: error was ", err)
	} else {
		fmt.Println("owner for dex.bonfida ", res2.NftOwner.ToBase58())
	}
	fmt.Printf("\n\n")

	/******************************************
	* IPFS.bonfida
	******************************************/

	if res, err = GetDomainKeySync("IPFS.bonfida", RecordVersion2); err != nil {
		fmt.Println("error was ", err)
	} else {
		fmt.Printf("pubkey for IPFS.bonand struct values are %+v\nfida ", res.PubKey.ToBase58())
	}

	if reverseLookUPStr, err = ReverseLookup(rpcClient, res.PubKey); err != nil {
		fmt.Println("reverseLookup: err was ", err)
	} else {
		fmt.Printf("reverse lookup for IPFS.bonfida pubkey(%s) is = `%s` ", res.PubKey.ToBase58(), reverseLookUPStr)
	}

	if res2, err = nm.Retrieve(rpcClient, res.PubKey); err != nil {
		fmt.Println("Retrieve: error was ", err)
	} else {
		fmt.Println("owner for IPFS.bonfida ", res2.NftOwner.ToBase58())
	}
	fmt.Printf("\n\n")
}

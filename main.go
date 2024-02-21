package main

import (
	"fmt"

	"github.com/blocto/solana-go-sdk/rpc"
)

func main() {
	var (
		res DomainKeyResult
		// publicKey common.PublicKey
		res2 RetrieveResult
		err  error
	)
	rpcClient := createConnection(rpc.MainnetRPCEndpoint)

	if res, err = GetDomainKeySync("bonfida", RecordVersion2); err != nil {
		fmt.Println("error was ", err)
	}
	fmt.Println("pubkey for bonfida ", res.PubKey.ToBase58())

	var nm NameRegistryState

	if res2, err = nm.Retrieve(rpcClient, res.PubKey); err != nil {
		fmt.Println("error was ", err)
	}
	fmt.Println("owner for bonfida ", res2.NftOwner.ToBase58())

	if res, err = GetDomainKeySync("dex.bonfida", RecordVersion2); err != nil {
		fmt.Println("error was ", err)
	}
	fmt.Println("pubkey for dex.bonfida ", res.PubKey.ToBase58())

	if res2, err = nm.Retrieve(rpcClient, res.PubKey); err != nil {
		fmt.Println("error was ", err)
	}
	fmt.Println("owner for dex.bonfida ", res2.NftOwner.ToBase58())

}

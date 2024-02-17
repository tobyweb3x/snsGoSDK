package main

import (
	"fmt"
	"log"
)

func main() {
	// client := createConnection(rpc.MainnetRPCEndpoint)

	s, err := GetDomainKeySync("bonfida", RecordVersion2)
	if err != nil {
		log.Fatal("error was: ", err)
	}
	fmt.Println("pubkey was ", s.DeriveResult.PubKey.String())
}

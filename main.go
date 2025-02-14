package main

import (
	"snsGoSDK/resolve"

	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	conn := rpc.New("https://mainnet.helius-rpc.com/?api-key=13af3657-7609-4ede-9305-6ea6c7a2243f")

	a, err := resolve.Resolve(conn, "sns-ip-5-wallet-12", resolve.ResolveConfig{})
	if err != nil {
		panic(err)
	}
	println(a.String())
}

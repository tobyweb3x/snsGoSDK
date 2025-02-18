package main

import (
	"fmt"
	"os"
	utils "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	conn := rpc.New(os.Getenv("RPC_ENDPOINT"))

	// a, err := resolve.Resolve(conn, "sns-ip-5-wallet-12", resolve.ResolveConfig{})
	// if err != nil {
	// 	panic(err)
	// }
	// println(a.String())

	// b, err := utils.ReverseLookup(conn,
	// 	solana.MustPublicKeyFromBase58("BxfSXLfrj3DsFVk6Pnqnt2b7F5AwoaSPtsbDnRCXfwbe"),
	// 	solana.PublicKey{},
	// )
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	println(b)
	// }

	c, err := utils.ReverseLookUpBatch(conn,
		[]solana.PublicKey{
			solana.MustPublicKeyFromBase58("Crf8hzfthWGbGbLTVCiqRqV5MVnbpHB1L9KQMd6gsinb"),
			solana.MustPublicKeyFromBase58("HPjEbJoeS77Qq31tWuS8pZCsY2yHAW2PcpAKBfETuwLa"),
			solana.MustPublicKeyFromBase58("JCqTzrANia2yfS5jDwpM76rFtyVvj4zu2nozVDk29wTh"),
			solana.MustPublicKeyFromBase58("2uSQkZRtJDYmBEbSg2WwMeWs2y21PNgGsUNoVRLDGRXZ"),
			solana.MustPublicKeyFromBase58("54obixuvJKGeJ6zFwYy1zb55G5c5z3B65MRXcc7fmaVU"),
			solana.MustPublicKeyFromBase58("5monfqudwcjVztfNa4nAyL4AwwymKgBspdR3RcvhKX4w"),
		},
	)
	if err != nil {
		fmt.Println("err->", err)
	} else {
		fmt.Println(c)
	}

}

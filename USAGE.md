## USAGE
This Go SDK mirrors the naming conventions and structure of the JS/TS SDK, so if you’re familiar with that, adapting here should feel intuitive. Always check for errors in usage, and if you’re unsure how a function behaves, refer to the corresponding test in the testing package for guidance.

> ### Resolve a domain owner

```go
import (
	"snsGoSDK/resolve"

    "github.com/gagliardetto/solana-go/rpc"
)
    conn := rpc.New(rpc.MainNetBeta_RPC)
	defer conn.Close()

    owner, err := resolve.Resolve(conn, "tobytobias", resolve.ResolveConfig{})
    // see the test for `Resolve()` for customization with the config.
```

> ### Get the Registry data for a domain

```go
import (
	"snsGoSDK/types"
	"snsGoSDK/utils"

    "github.com/gagliardetto/solana-go/rpc"
)
    conn := rpc.New(rpc.MainNetBeta_RPC)
	defer conn.Close()

    // utils.GetDomainKeySync("mrwick", types.V1)
    // utils.GetDomainKeySync("mrwick.sol", types.V2)
    domainKey, err := utils.GetDomainKeySync("mrwick", types.V0) // types.VO is the equivalent for `undefined` in the TS SDK
    nm := spl.NameRegistryState{}
    registry, err := nm.Retrieve(conn, domainKey.PubKey)
    nameRegState, nftOwner := registry.Registry, registry.NftOwner
```

> ### Get the Registry data for a domain (batch)

```go
import (
	"snsGoSDK/types"
	"snsGoSDK/utils"

    "github.com/gagliardetto/solana-go/rpc"
)
    conn := rpc.New(rpc.MainNetBeta_RPC)
	defer conn.Close()

    nm := spl.NameRegistryState{}
    registries, err := nm.RetrieveBatch(conn, []solana.PublicKey{})

    Note: spl.NameRegistryState.RetrieveBatch returns a slice of pointer to `NameRegistryState` (order of `[]solana.PublicKey{}` in the param is preserved).
```

> ### Reverse Lookup

```go
import (
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)
    conn := rpc.New(rpc.MainNetBeta_RPC)
	defer conn.Close()

	domainKey := solana.MustPublicKeyFromBase58("Crf8hzfthWGbGbLTVCiqRqV5MVnbpHB1L9KQMd6gsinb")
	domainName, err := utils.ReverseLookup(conn, domainKey, solana.PublicKey{}) // bonfida
```

> ### Reverse Lookup (batch)

```go
import (
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)
    conn := rpc.New(rpc.MainNetBeta_RPC)
	defer conn.Close()

	domainNames, err := utils.ReverseLookUpBatch(conn, []solana.PublicKey{
		solana.MustPublicKeyFromBase58("Crf8hzfthWGbGbLTVCiqRqV5MVnbpHB1L9KQMd6gsinb"),
		solana.MustPublicKeyFromBase58("HPjEbJoeS77Qq31tWuS8pZCsY2yHAW2PcpAKBfETuwLa"),
		solana.MustPublicKeyFromBase58("JCqTzrANia2yfS5jDwpM76rFtyVvj4zu2nozVDk29wTh"),
		solana.MustPublicKeyFromBase58("2uSQkZRtJDYmBEbSg2WwMeWs2y21PNgGsUNoVRLDGRXZ"),
		solana.MustPublicKeyFromBase58("54obixuvJKGeJ6zFwYy1zb55G5c5z3B65MRXcc7fmaVU"),
		solana.MustPublicKeyFromBase58("45LYEaK4ZwBiymEXNNXPyPKUh8yxaDTVutJbXYwooHLE"),
	})
	// domainNames = []string{
	// 	"bonfida",
	// 	"tobytobias",
	// 	"menbehindwoman",
	// 	"grimmest",
	// 	"niftydegen",
	// 	"mrwick",
	// }

    also order is preserved.
```

> ### Get all domains of a user

```go
import (
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)
    conn := rpc.New(rpc.MainNetBeta_RPC)
	defer conn.Close()

	domainKeyAndReverse, err := utils.GetDomainKeysWithReverses(conn, solana.MustPublicKeyFromBase58("..."))
	//  domainKeyAndReverse is of type
	//  type GetDomainKeysWithReversesResult struct {
	// 	DomainKey solana.PublicKey
	// 	Domain string
	// }
```

> ### Get single & multiple Favorite domain

```go
import (
	"snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)
    conn := rpc.New(rpc.MainNetBeta_RPC)
	defer conn.Close()

	favDomain, err := utils.GetFavoriteDoamin(conn, solana.MustPublicKeyFromBase58("..."))
	favDomains, err := utils.GetMultipleFavoriteDomain(conn, []solana.PublicKey{})

    also order is preserved.
```

> ### Get all subdomains

```go
import (
    "snsGoSDK/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)
    conn := rpc.New(rpc.MainNetBeta_RPC)
	defer conn.Close()

	domainKeyAcc, err := utils.GetDomainKeySync("tobytobias", types.V0)
	subdomains, err := utils.FindSubdomains(conn, domainKeyAcc.PubKey)
    // subdomains is of type []string
```

> ### Burn a domain

```go
import (
	"snsGoSDK/bindings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)
    conn := rpc.New(rpc.MainNetBeta_RPC)
	defer conn.Close()

	owner, burnDst := solana.MustPublicKeyFromBase58("HKKp49qGWXd639QsuH7JiLijfVW5UtCVY4s1n2HANwEA"),
			solana.MustPublicKeyFromBase58("3Wnd5Df69KitZfUoPYZU438eFRNwGHkhLnSAWL65PxJX")
		ix, err := bindings.BurnDomain(
			"bonfida",
			owner,
			burnDst,
		)

		recent, err := conn.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)

		tx, err := solana.NewTransactionBuilder().
			AddInstruction(ix).
			SetRecentBlockHash(recent.Value.Blockhash).
			SetFeePayer(owner).Build()


		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				p := solana.MustPrivateKeyFromBase58(os.Getenv("TEST_PRIVATE_KEY"))
				return &p
			},
		)

		simTxn, err := conn.SimulateTransactionWithOpts(
			context.TODO(),
			tx,
			&rpc.SimulateTransactionOpts{},
		)
```

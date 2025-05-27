This SDK follows the same conventions in directory naming, function names and variable with the JS/TS SDK, so it should not be very hard jumping on how things are don here;

> #### Get the Registry data for a domain

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

> #### Get the Registry data for a domain (batch)

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

    Note: spl.NameRegistryState.RetrieveBatch can be used to retrieve multiple name registries at once. Pass the connection, and an slice of domain name public keys as arguments to the function, it returns a slice of pointer to `NameRegistryState` (order of `[]solana.PublicKey{}` in the param is preserved).
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

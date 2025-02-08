![](Screenshot%20from%202024-02-26%2011-40-12.png)
![](Screenshot%20from%202024-02-26%2011-40-45.png)

I have accounted for all of that 👆, for now.

To test run;

> go build && ./snsGOSDK

The is the Js test of the library to corelate results.

```typescript
import { clusterApiUrl, Connection, PublicKey } from "@solana/web3.js";
// @ts-ignore
import {
  getDomainKeySync,
  NameRegistryState,
  RecordVersion,
  reverseLookup,
  getHandleAndRegistryKey,
  getTwitterRegistry,
  getAllDomains,
  getDomainKeysWithReverses,
} from "@bonfida/spl-name-service";

async function createConnection() {
  return new Connection(clusterApiUrl("mainnet-beta"));
}

const testOthers = async (conn: Connection) => {
  const rangeValues = [
    "bonfida.sol",
    "solana.sol",
    "01.sol",
    "dex.solana.sol",
    "dex.bonfida.sol",
    "wallet-guide-5.sol",
    "sub-0.wallet-guide-3.sol",
  ];
  for (const [index, domain] of rangeValues.entries()) {
    console.log(`${index + 1}). domain - ${domain}`);

    console.log("*****GetDomainKeySync*****");
    const { pubkey } = getDomainKeySync(domain, RecordVersion.V2);
    console.log(`pubkey for "${domain} is"--- "${pubkey.toBase58()}"`);

    try {
      console.log("*****ReverseLookUp*****");
      const reverse = await reverseLookup(conn, pubkey);
      console.log(
        `\nreverselookup for PublicKey --${pubkey} is = "${reverse}"`
      );
    } catch (error) {
      console.log(`Error from domain reverseLookup----: \n${error}`);
      continue;
    }

    try {
      console.log("*****NameStateRegistry.Retrieve*****");
      const { registry, nftOwner } = await NameRegistryState.retrieve(
        conn,
        pubkey
      );
      console.log(
        `Owner's PublicKey for domain ${domain} is "${registry.owner.toBase58()}"`
      );
      if (nftOwner) {
        console.log("nftOwner----- ", nftOwner.toBase58(), "\n");
      }
    } catch (error) {
      console.log(
        `Error from domain NameRegistryState.retrieve----: \n${error}`
      );
      continue;
    }

    try {
      console.log("*****GetAllDomains*****");
      const { registry } = await NameRegistryState.retrieve(conn, pubkey);
      const domains = await getAllDomains(conn, registry.owner);
      console.log("All domains of ${registry.owner} ----->");
      for (const domain of domains) {
        console.log(`${domain.toBase58()}`);
      }
    } catch (error) {
      console.log(`Error from domain getAllDomains----: \n${error}`);
      continue;
    }

    try {
      console.log("*****getDomainKeysWithReverses*****");
      const { registry, nftOwner } = await NameRegistryState.retrieve(
        conn,
        pubkey
      );
      const domainsWithReverses = await getDomainKeysWithReverses(
        conn,
        registry.owner
      );
      console.log("All domains of ${registry.owner} ----->");

      for (const domain of domainsWithReverses) {
        console.log(`${domain.pubKey} "--**********--"  ${domain.domain}`);
      }
    } catch (error) {
      console.log(
        `Error from domain getDomainKeysWithReverses----: \n${error}`
      );
      continue;
    }
  }
};

const testTwitterResolving = async (conn: Connection) => {
  const rangeValues = ["oluwatobialone"];
  for (const domain of rangeValues) {
    console.log("\n");
    const reg = await getTwitterRegistry(conn, domain);
    console.log(
      `public key associated to the Twitter handle @${domain}----`,
      reg.owner.toBase58()
    );

    const [handle] = await getHandleAndRegistryKey(conn, reg.owner);
    console.log(
      `Twitter handle associated to the public key ${domain} is----`,
      handle,
      "\n\n"
    );
  }
};

(async () => {
  const conn = await createConnection();
  await testOthers(conn);
  await testTwitterResolving(conn);
})();
```

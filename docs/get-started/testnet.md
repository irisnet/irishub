---
order: 4
---

# Join The Testnet

We have 2 testnets: Fuxi and Nyancat.

Since the launch of mainnet, **Fuxi Testnet** starts to operate as a stable application testnet which has the same version as the mainnet, so that the service providers of IRISnet can develop their apps on or interact with IRIShub without running a node or lcd instance.

However there is also a need for validators to test the new version of IRIShub before it can be relased to production, and this is **Nyancat Testnet**'s focus. And new validators can also use the Nyancat Testnet to practice the validator operations.

## Install

We use different bech32 prefixes to distinguish the mainnet and testnet, all you need to do is to run the following command in the [irishub](https://github.com/irisnet/irishub) source root before [building or installing](install.md) the iris binaries:

```bash
source scripts/setTestEnv.sh # to build or install the testnet version
```

## Fuxi Testnet

There are no options to run nodes to connect to the Fuxi Testnet, you can use the public RPC and LCD to develop and test your apps.

- RPC: <http://rpc.testnet.irisnet.org:80>

- LCD: <https://lcd.testnet.irisnet.org/swagger-ui/>

- Explorer: <https://testnet.irisplorer.io>

## Nyancat Testnet

Please refer to [How to join nyancat testnet](https://github.com/irisnet/testnets/tree/master/nyancat#how-to-join-nyancat-testnet)

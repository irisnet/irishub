---
order: 2
---

# Join The Mainnet

:::tip
We assume that you already have `iris` installed, or you need to [install iris](../software/How-to-install-irishub.md) first.
:::

## Run a Full Node

```bash
# initialize node configurations 
iris init --moniker=<your_custom_name> --chain-id=irishub

# download mainnet public config.toml and genesis.json
curl -o ~/.iris/config/config.toml https://raw.githubusercontent.com/irisnet/betanet/master/config/config.toml
curl -o ~/.iris/config/genesis.json https://raw.githubusercontent.com/irisnet/betanet/master/config/genesis.json

# start the node (you can also use "nohup" to run in the background)
iris start
```

:::tip
You may see some connection errors, it does not matter, the P2P network is trying to find available connections

[Advanced Configurations](#TODO)
:::

## Upgrade to Validator Node

You now have an active full node. What's the next step? 

If you have participated in the genesis file generation process, you should be a validator once you are fully synced. 

If you miss the genesis file generation process, you can still upgrade your full node to become an IRISnet Validator. The top 100 validators have the ability to propose new blocks to the IRIS Hub. 

Please follow this [instruction](Validator-Node.md) to upgrade your full node to validator node.

## Deploy IRIShub Monitor

Please follow this [guide](../software/monitor.md) to deploy IRIHub Monitor.

## Setup a Sentry Node

A validator is under the risk of being attacked. You could follow this [guide](../software/sentry.md) to setup a sentry node to protect yourself.

## Use a KMS
If you plan to use a KMS (key management system), you should go through these steps first: [Using a KMS](../software/kms/kms.md).

##  Useful Links

- Riot chat: #irisvalidators:matrix.org

- Explorer: <https://www.irisplorer.io>
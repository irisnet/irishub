---
order: 3
---

# Join The Mainnet

:::tip
**Requirements:** [install iris](install.md)
:::

## Run a Full Node

```bash
# initialize node configurations
iris init --moniker=<your-custom-name> --chain-id=irishub

# download mainnet public config.toml and genesis.json
curl -o ~/.iris/config/config.toml https://raw.githubusercontent.com/irisnet/mainnet/master/config/config.toml
curl -o ~/.iris/config/genesis.json https://raw.githubusercontent.com/irisnet/mainnet/master/config/genesis.json

# start the node (you can also use "nohup" or "systemd" to run in the background)
iris start
```

:::tip
You may see some connection errors, it does not matter, the P2P network is trying to find available connections

Try to add some of the [Community Peers](https://github.com/irisnet/mainnet/blob/master/config/community-peers.md) to `persistent_peers` in the config.toml
:::

:::tip
It will take a long time to sync from scratch, you can also download the [mainnet data snapshot](#TODO) to reduce the time spent on synchronization
:::

## Upgrade to Validator Node

### Create a Wallet

You can [create a new wallet](../cli-client/keys.md#create-a-new-key) or [import an existing one](../cli-client/keys.md#recover-an-existing-key-from-seed-phrase), then get some IRIS from the exchanges or anywhere else into the wallet you just created, .e.g.

```bash
# create a new wallet
iriscli keys add <key-name>
```

:::warning
**Important**

write the seed phrase in a safe place! It is the only way to recover your account if you ever forget your password.
:::

### Confirm your node has caught-up

```bash
# if you have not installed jq
# apt-get update && apt-get install -y jq

# if the output is false, means your node has caught-up
iriscli status | jq .sync_info.catching_up
```

### Create Validator

Only if your node has caught-up, you can run the following command to upgrade your node to be a validator.

```bash
iriscli stake create-validator \
    --pubkey=$(iris tendermint show-validator) \
    --moniker=<your-validator-name> \
    --amount=<amount-to-be-delegated, e.g. 10000iris> \
    --commission-rate=0.1 \
    --gas=100000 \
    --fee=0.6iris \
    --chain-id=irishub \
    --from=<key-name> \
    --commit
```

:::warning
**Important**

Backup the `config` directory located in your iris home (default ~/.iris/) carefully! It is the only way to recover your validator.
:::

If there are no errors, then your node is now a validator or candidate (depending on whether your delegation amount is in the top 100)

Read more:

- Concepts
  - [General Concepts](../concepts/general-concepts.md)
  - [Validator FAQ](../concepts/validator-faq.md)
- Validator Security
  - [Sentry Nodes (DDOS Protection)](../concepts/sentry-nodes.md)
  - [Key Management](../tools/kms.md)

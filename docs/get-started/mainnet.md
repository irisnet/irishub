---
order: 3
---

# Join The Mainnet

:::tip
**Requirements:** [install iris](install.md)
:::

## Run a Full Node

### Start node from genesis

:::tip
You must use irishub [v1.0.1](https://github.com/irisnet/irishub/releases/tag/v1.0.1) to initialize your node.
:::

```bash
# initialize node configurations
iris init <moniker> --chain-id=irishub-1

# download mainnet public config.toml and genesis.json
curl -o ~/.iris/config/config.toml https://raw.githubusercontent.com/irisnet/mainnet/master/config/config.toml
curl -o ~/.iris/config/genesis.json https://raw.githubusercontent.com/irisnet/mainnet/master/config/genesis.json

# start the node (you can also use "nohup" or "systemd" to run in the background)
iris start
```

Next, your node will process all chain upgrades. Between each upgrade, you must use the specified version to catch up with the block. Don't worry about using the old version at the upgrade height, the node will be halted automatically.

| Proposal | Start height | Upgrade height | irishub version |
| -------- | ------------ | -------------- | ----- |
| genesis  |  9146456     |  9593205  | [v1.0.1](https://github.com/irisnet/irishub/releases/tag/v1.0.1) |
| [#1](https://irishub.iobscan.io/#/ProposalsDetail/1)  |  9593206     |    | [v1.1.0](https://github.com/irisnet/irishub/releases/tag/v1.1.0),[v1.1.1](https://github.com/irisnet/irishub/releases/tag/v1.1.1)|

:::tip
You may see some connection errors, it does not matter, the P2P network is trying to find available connections

Try to add some of the [Community Peers](https://github.com/irisnet/mainnet/blob/master/config/community-peers.md) to `persistent_peers` in the config.toml
:::

## Upgrade to Validator Node

### Create a Wallet

You can [create a new wallet](../cli-client/keys.md#create-a-new-key) or [import an existing one](../cli-client/keys.md#recover-an-existing-key-from-seed-phrase), then get some IRIS from the exchanges or anywhere else into the wallet you just created, .e.g.

```bash
# create a new wallet
iris keys add <key-name>
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
iris status | jq .sync_info.catching_up
```

### Create Validator

Only if your node has caught-up, you can run the following command to upgrade your node to be a validator.

```bash
iris tx staking create-validator \
    --pubkey=$(iris tendermint show-validator) \
    --moniker=<your-validator-name> \
    --amount=<amount-to-be-delegated, e.g. 10000iris> \
    --min-self-delegation=1 \
    --commission-max-change-rate=0.1 \
    --commission-max-rate=0.1 \
    --commission-rate=0.1 \
    --gas=100000 \
    --fees=0.6iris \
    --chain-id=irishub-1 \
    --from=<key-name>
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
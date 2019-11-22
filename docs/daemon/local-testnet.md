---
order: 3
---

# Local Testnet

For testing or developing purpose, you may want to setup a local testnet.

## Single Node Testnet

**Requirements:**

- [Install iris](../get-started/install.md)

:::tip
We use the default [home directory](intro.md#home-directory) for all the following examples
:::

### iris init

Initialize the genesis.json file that will help you to bootstrap the network

```bash
iris init --chain-id=testing --moniker=testing
```

### create a key

Create a key to hold your validator account

```bash
iriscli keys add MyValidator
```

### iris add-genesis-account

Add that key into the genesis.app_state.accounts array in the genesis file

:::tip
this command lets you set the number of coins. Make sure this account has some iris which is the only staking coin on IRISnet
:::

```bash
iris add-genesis-account $(iriscli keys show MyValidator --address) 100000000iris
```

### iris gentx

Generate the transaction that creates your validator. The gentxs are stored in `~/.iris/config/gentx/`

```bash
iris gentx --name MyValidator
```

### iris collect-gentxs

Add the generated staking transactions to the genesis file

```bash
iris collect-gentxs
```

### iris start

Now it‘s ready to start `iris`

```bash
iris start
```

### iris unsafe-reset-all

You can use this command to reset your node, including the local blockchain database, address book file, and resets priv_validator.json to the genesis state.

This is useful when your local blockchain database somehow breaks and you are not able to sync or participate in the consensus.

```bash
iris unsafe-reset-all
```

### iris reset

Unlike [iris unsafe-reset-all](#iris-unsafe-reset-all), this command allows you to reset the blockchain state of your node to a specified height, so you can fix your blockchain database much faster.

```bash
# e.g. reset the blockchain state to height 100
iris reset --height 100
```

And there is another option to fix the blockchain database, if you got a `Wrong Block.Header.AppHash` error on the Mainnet, confirm you are using the correct [Mainnet Version](../get-started/install.md#latest-version), then restart your node by:

```bash
iris start --replay-last-block
```

### iris tendermint

Query the unique node id which can be used in p2p connection, e.g. the `seeds` and `persistent_peers` in the [config.toml](intro.md#cnofig-toml) are formatted as `<node-id>@ip:26656`.

The node id is stored in the [node_key.json](intro.md#node_key-json).

```bash
iris tendermint show-node-id
```

Query the [Tendermint Pubkey](../concepts/validator-faq.md#tendermint-key) which is used to [identify your validator](../cli-client/stake/create-validator.md), and the corresponding private key will be used to sign the Pre-vote/Pre-commit in the consensus.

The [Tendermint Key](../concepts/validator-faq.md#tendermint-key) is stored in the [priv_validator.json](intro.md#priv_validator-json) which is [required to be backed up](../concepts/validator-faq.md#how-to-backup-the-validator) once you become a validator.

```bash
iris tendermint show-validator
```

Query the bech32 prefixed validator address

```bash
iris tendermint show-address
```

### iris export

Please refer to [Export Blockchain State](export.md)

## Multiple Nodes Testnet

**Requirements:**

- [Install iris](../get-started/install.md)
- [Install jq](https://stedolan.github.io/jq/download/)
- [Install docker](https://docs.docker.com/engine/installation/)
- [Install docker-compose](https://docs.docker.com/compose/install/)

### Build and Init

```bash
# Work from the irishub repo
cd $GOPATH/src/github.com/irisnet/irishub

# Build the linux binary in ./build
make build_linux

# Quick init a 4-node testnet configs
make testnet_init
```

The `make testnet_init` generates config files for a 4-node testnet in the `./build/nodecluster` directory by calling the `iris testnet` command:

```bash
$ tree -L 3 build/nodecluster/
build/nodecluster/
├── gentxs
│   ├── node0.json
│   ├── node1.json
│   ├── node2.json
│   └── node3.json
├── node0
│   ├── iris
│   │   ├── config
│   │   └── data
│   └── iriscli
│       ├── key_seed.json
│       └── keys
├── node1
│   ├── iris
│   │   ├── config
│   │   └── data
│   └── iriscli
│       └── key_seed.json
├── node2
│   ├── iris
│   │   ├── config
│   │   └── data
│   └── iriscli
│       └── key_seed.json
└── node3
    ├── iris
    │   ├── config
    │   └── data
    └── iriscli
        └── key_seed.json
```

### Start

```bash
make testnet_start
```

This command creates a 4-node network using the ubuntu:16.04 docker image. The ports for each node are found in this table:

| Node      | P2P Port | RPC Port |
| --------- | -------- | -------- |
| irisnode0 | 26656    | 26657    |
| irisnode1 | 26659    | 26660    |
| irisnode2 | 26661    | 26662    |
| irisnode3 | 26663    | 26664    |

To update the binary, just rebuild it and restart the nodes:

```bash
make build_linux testnet_start
```

### Stop

To stop all the running nodes:

```bash
make testnet_stop
```

### Reset

To stop all the running nodes and reset the network to the genesis state:

```bash
make testnet_unsafe_reset
```

### Clean

To stop all the running nodes and delete all the files in the `build/` directory:

```bash
make testnet_clean
```

# Local Testnet

For testing or developing purpose, you may want to setup a local testnet.

## Single Node Testnet

**Requirements:**

- [Install iris](install.md)

```bash
# Initialize the genesis.json file that will help you to bootstrap the network
iris init --chain-id=testing --moniker=testing

# Create a key to hold your validator account
iriscli keys add MyValidator

# Add that key into the genesis.app_state.accounts array in the genesis file
# NOTE: this command lets you set the number of coins. Make sure this account has some iris
# which is the only staking coin on IRISnet
iris add-genesis-account $(iriscli keys show MyValidator --address) 100000000iris

# Generate the transaction that creates your validator
# The gentxs are stored in ~/.iris/config/gentx/
iris gentx --name MyValidator

# Add the generated bonding transactions to the genesis file
iris collect-gentxs

# Now it‘s ready to start `iris`
iris start
```

## Multiple Nodes Testnet

**Requirements:**

- [Install iris](install.md)
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

The `make testnet_init` creates files for a 4-node testnet in `./build` by calling the `iris testnet` command. This outputs a handful of files in the `./build/nodecluster` directory:

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

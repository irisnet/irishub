# IRIS Daemon

## Introduction

The iris executable is the entry point for running a IRISnet network node. All the validator nodes and full nodes need to install the iris command and launching the daemon to join the IRISnet network. You can also use this command to start your own test network locally. If you need to join the IRISnet testnet, please refer to [get-started](../get-started/README.md).

## How to start an IRISnet network locally

### Initialize node

First you need to create a account as the corresponding validator operator for yourself.
```bash
iriscli keys add {account_name}
```
You can get the account information, including account address, public key address and mnemonic
```
NAME:	TYPE:	ADDRESS:						PUBKEY:
account_name	local	faa13t6jugwm5uu3h835s5d4zggkklz6rpns59keju	fap1addwnpepqdne60eyssj2plrsusd8049cs5hhhl5alcxv2xu0xmzlhphy9lyd5kpsyzu
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

witness exotic fantasy gaze brass zebra adapt guess drip quote space payment farm argue pear actress garage smile hawk bid bag screen wonder person
```

Initialize the configuration files such as genesis.json and config.toml
```bash
iris init --home={path_to_your_home} --chain-id={your_chain_id} --moniker={your_node_name}
```
This command will create the corresponding files in the home directory.

Create the `CreateValidator` transaction and sign the transaction by the validator operator account you just created
```bash
iris gentx --name={account_name} --home={path_to_your_home}
```
This commond will generate the transaction in the directory：{path_to_your_home}/config/gentx

### Config genesis

Use the following command to modify the genesis.json file to assign the initial account balance to the above validator operator account, such as: 150 iris
```bash
iris add-genesis-account faa13t6jugwm5uu3h835s5d4zggkklz6rpns59keju 150iris
```

```json
    {
      "accounts": [
        {
          "address": "faa13t6jugwm5uu3h835s5d4zggkklz6rpns59keju",
          "coins": ["150iris"],
          "sequence_number": "0",
          "account_number": "0"
        }
      ]
    }
```

Configuring validator information
```bash
iris collect-gentxs --home={path_to_your_home}
```
This command reads the `CreateValidator` transaction under folder {path_to_your_home}/config/gentx and writes it to genesis.json to complete the assignment of the initial validators in the Genesis block.

### Boot node

After completing the above configuration, start the node with the following command
```bash
iris start --home {path_to_your_home}
```
After the command is executed, iris generate the genesis block according to the genesis.json configured under the home and update the application state. Then, according to the current block height of the chain pointed by the chain-id, it either starts to synchronize with other peers, or enter the process of produce the first block.

## Start a multi-node network

To start a multi-node IRISnet network locally, you need to follow the steps below：

* Prepare home directory: create a unique home directory for each node
* Initialization: according to the above `Initialize node` steps, initialize each node in the respective home directory (note that you need to use the same chain-id)
* Configure genesis: select the home directory of one of the nodes to configure the genesis, refer to the `Config genesis` steps, configure the account address and initial balance of each node, and copy the files from each node {path_to_your_home}/config/gentx to the current one. Execute `iris collect-gentxs` to generate the final genesis.json. Finally, copy the genesis.json to the home directory of each node, overwriting the original genesis.json.
* Configure config.toml: modify {path_to_your_home}/config/config.toml in the home directory of each node, assign different ports to each node, query the node-id of each node through `iris tendermint show-node-id`, and then add the `node-id@ip:port` of other nodes in `persistent_peers` so that nodes can connect to each other.

In order to simplify the above configuration process, you can use the following commands to automate the initialization of local multi-node and the configuration of genesis and config.toml:
```bash
iris testnet --v 4 --output-dir ./output --chain-id irishub-test --starting-ip-address 127.0.0.1
```

To start the node to join the public IRIS Testnet, please refer to [Full-Node](../get-started/Full-Node.md)

## Introduction of home directory 

The home directory is the working directory of the iris node. The home directory contains all the configuration information and all the data that the node runs.

In the iris command, you can specify the home directory of the node by using flag `--home`. If you run multiple nodes on the same machine, you need to specify a different home directory for them. If the `--home` flag is not specified in the iris command, the default value `$HOME/.iris` is used as the home directory used by this iris command.

The `iris init` command is responsible for initializing the specified `--home` directory and creating the default configuration files. Except the `iris init` command, the home directory used by any other `iris` sub commands must be initialized, otherwise an error will be reported.

The data of the iris node is stored in the `data` directory of the home, including blockchain data, application layer data, and index data. All configuration files are stored in the home/`config` directory:

### genesis.json

genesis.json defines the genesis block data, which specifies the system parameters such as chain_id, consensus parameters, initial account token allocation, creation of validators, and parameters for stake/slashing/gov/upgrade. See [genesis-file](../features/basic-concepts/genesis-file.md) for details.

### node_key.json

node_key.json is used to store the node's key. The node-id queried by `iris tendermint show-node-id` is derived by the key, which is used to indicate the unique identity of the node. It is used in p2p connection.

### pri_validator.json

pri_validator.json is the key that the validator will use to sign Pre-vote/Pre-commit in each round of consensus voting. As the consensus progresses, the tendermint consensus engine will continuously update `last_height`/`last_round`/ `last_step` and other data.

### config.toml

config.toml is the non-consensus configuration information of the node. Different nodes can configure themselves according to their own situation. Common modifications are `persistent_peers`/`moniker`/`laddr`.

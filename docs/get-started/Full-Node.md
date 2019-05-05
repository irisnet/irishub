# Setup A Full Node

Before setting up your validator node on IRIShub, make sure you already had `iris` installed by following this [guide](../software/How-to-install-irishub.md) 

## Init Your Node

These instructions are for setting up a brand new full node from scratch.

First, initialize the node and create the necessary config files:

```
iris init --moniker=<your_custom_name> --home=<iris_home> --chain-id=<chain_id>
```
> The chain-id of mainnet is `irishub`
> Note: Only ASCII characters are supported for the `--moniker`. Using Unicode characters will render your node unreachable.

The default <iris_home> is `~/.iris` , You can edit this `moniker` later in `~/.iris/config/config.toml`.

Your full node has been initialized!

## Get Configuration Files

After initializing your node, please download the genesis file and the config file to join in the IRISnet.

```
cd <iris_home>/config/
rm genesis.json
rm config.toml
wget https://raw.githubusercontent.com/irisnet/betanet/master/config/genesis.json
wget https://raw.githubusercontent.com/irisnet/betanet/master/config/config.toml
```
## Edit Your Config File

You could customize the `moniker` and `external_address` fields. 

```
# A custom human readable name for this node
moniker = "<your_custom_name>"
external_address = "<your_public_IP>:26656"
```

Optional:
Set `addr_book_strict` to `false` to make peering more easily.

```
addr_book_strict = false
```


### Add Seed Nodes

Your node needs to know how to find more peers. You'll need to add healthy seed nodes to `<iris_home>/config/config.toml`. Here are some seed nodes you can use:

```
6a6de770deaa4b8c061dffd82e9c7f4d40c2165d@seed-1.mainnet.irisnet.org:26656
a17d7923293203c64ba75723db4d5f28e642f469@seed-2.mainnet.irisnet.org:26656
```

Meanwhile, you could add some known full node as `Persistent Peer`. Your node could connect to `sentry node` as `persistent peers`.


###  Enable Port

You will need to set `26656` port to get connected with other peers and `26657` to query information of Tendermint.

## Run a Full Node

Start the full node with this command:

```
iris start --home=<iris_home> > iris.log &
```

Check that everything is running smoothly:

```
iriscli status
```
You could see the following
```json
{
  "node_info": {
    "protocol_version": {
      "p2p": "5",
      "block": "8",
      "app": "0"
    },
    "id": "8fa36b85e98f986b70889da52b733fa925908947",
    "listen_addr": "tcp://0.0.0.0:26656",
    "network": "irishub",
    "version": "0.27.3",
    "channels": "4020212223303800",
    "moniker": "test",
    "other": {
      "tx_index": "on",
      "rpc_address": "tcp://0.0.0.0:26657"
    }
  },
  "sync_info": {
    "latest_block_hash": "DF2F64D56863C5516586112B9A954DFB2257C65FF178267E75D85D160E5E0E2B",
    "latest_app_hash": "",
    "latest_block_height": "1",
    "latest_block_time": "2019-01-23T03:42:17.268038878Z",
    "catching_up": false
  },
  "validator_info": {
    "address": "7B9052D259643E5B9AF0BD481B843C89B27AACAA",
    "pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "Mt9WvNPDd8F4Dcq7lP+GFIhW0/K4jAt8nTq/ljut94E="
    },
    "voting_power": "100"
  }
}
```
If you see the 	`catching_up` is `false`, it means your node is fully synced with the network, otherwise your node is still downloading blocks. Once fully synced, you could upgrade your node to a validator node. The instructions is in [here](Validator-Node.md).	

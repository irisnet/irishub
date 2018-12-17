# Setup A Full-node

Before setting up your validator node, make sure you already had **Iris** installed by following this [guide](Install-the-Software.md)

## Init Your Node

These instructions are for setting up a brand new full node from scratch.

First, initialize the node and create the necessary config files:

```
iris init --moniker=<your_custom_name> --home=$IRISHOME --chain-id=<chain-id>
```

> Note: Only ASCII characters are supported for the `--moniker`. Using Unicode characters will render your node unreachable.

The default \$IRISHOME is `~/.iris` , You can edit this `name` later, in the `~/.iris/config/config.toml` file:

Your full node has been initialized!

## Get Configuration Files


After intializing your node, please download the genesis file and the config file to join in the testnet.

```
cd $IRISHOME/config/
rm genesis.json
rm config.toml
wget https://raw.githubusercontent.com/irisnet/testnets/master/fuxi/fuxi-6000/config/config.toml
wget https://raw.githubusercontent.com/irisnet/testnets/master/fuxi/fuxi-6000/config/genesis.json
```
## Edit Your Config File

You could customized the `moniker` and `external_address` fields. 

```
# A custom human readable name for this node
moniker = "<your_custom_name>"
external_address = "your-public-IP:26656"
```


Optional:
Set `addr_book_strict` to `false` to make peering more easily.

```
addr_book_strict = false
```


### Add Seed Nodes

Your node needs to know how to find more peers. You'll need to add healthy seed nodes to `$IRISHOME/config/config.toml`. Here are some seed nodes you can use:

```
c16700520a810b270206d59f0f02ea9abd85a4fe@ts-1.bianjie.ai:26656
a12cfb2f535210ea12731f94a76b691832056156@ts-2.bianjie.ai:26656
```

Meanwhile, you could add some known full node as `Persistent Peer`. Your node could connect to `sentry node` as `persistent peers`.


###  Enable Port

You will need to set `26656` port to get connected with other peers and `26657` to query information of Tendermint.

## Run a Full Node

Start the full node with this command:

```
iris start --home=$IRISHOME > iris.log
```

Check that everything is running smoothly:

```
iriscli status
```
You could see the following
```
{"node_info":{"id":"1c40d19d695721fc3e3ce44cbc3f446f038b36e4","listen_addr":"172.31.0.190:26656","network":"fuxi-6000","version":"0.22.6","channels":"4020212223303800","moniker":"name","other":["amino_version=0.10.1","p2p_version=0.5.0","consensus_version=v1/0.2.2","rpc_version=0.7.0/3","tx_index=on","rpc_addr=tcp://0.0.0.0:46657"]},"sync_info":{"latest_block_hash":"41117D8CB54FA54EFD8DEAD81D6D83BDCE0E63AC","latest_app_hash":"95D82B8AC8B64C4CD6F85C1D91F999C2D1DA4F0A","latest_block_height":"1517","latest_block_time":"2018-09-07T05:44:27.810641328Z","catching_up":false},"validator_info":{"address":"3FCCECF1A27A9CEBD394F3A0C5253ADAA8392EB7","pub_key":{"type":"tendermint/PubKeyEd25519","value":"wZp1blOEwJu4UuqbEmivzjUMO1UwUK4C0jRH96HhV90="},"voting_power":"100"}}
```
If you see the 	`catching_up` is `false`, it means your node is fully synced with the network, otherwise your node is still downloading blocks. Once fully synced, you could upgrade your node to a validator node. The instructions is in [here](Validator-Node.md).	

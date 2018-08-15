# Setup A Full-node

Before setting up your validator node, make sure you've already installed  **Iris** by this [guide](https://github.com/irisnet/testnets/blob/master/testnets/docs/install%20iris.md)

### Step 2: Setting Up Your Node

These instructions are for setting up a brand new full node from scratch.

First, initialize the node and create the necessary config files:

```
iris init --name <your_custom_name> --home=<IRISHOME>
```

> Note: Only ASCII characters are supported for the `--name`. Using Unicode characters will render your node unreachable.

The default \$IRISHOME is `~/.iris` , You can edit this `name` later, in the `~/.iris/config/config.toml` file:

```
# A custom human readable name for this node
moniker = "<your_custom_name>"
external_address = "<your-public-IP>"
```


Optional:
Set `addr_book_strict` to `false` to make peering more easily.

```
addr_book_strict = false
```
Your full node has been initialized!

### Get Configuration Files

If you want to be part of the genesis file geneartion process, please follow this [guide](https://github.com/kidinamoto01/testnets-1/blob/master/testnets/docs/Genesis%20Generation%20Process.md) to sumbmit a json file. Otherwise, you could always send related transaction to become a validator later.

After the genesis file generation process is finished, please download the genesis and the default config file.

```
cd $IRISHOME/config/
rm genesis.json
rm config.toml
wget https://raw.githubusercontent.com/irisnet/testnets/master/testnets/fuxi-2000/config/config.toml
wget https://raw.githubusercontent.com/irisnet/testnets/master/testnets/fuxi-2000/config/genesis.json
```

### Add Seed Nodes

Your node needs to know how to find peers. You'll need to add healthy seed nodes to `$IRISHOME/config/config.toml`. Here are some seed nodes you can use:

```
c16700520a810b270206d59f0f02ea9abd85a4fe@35.165.232.141:26656
```

Meanwhile, you could add some known full node as `Persistent Peer`. Your node could connect to `sentry node` as `persistent peers`.


Chang the `external_address` to your `public IP:26656`.


### Run a Full Node

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
{"node_info":{"id":"71b188e9fdefd939453b3cd10c0eae45a8d02a2b","listen_addr":"172.31.0.190:26656","network":"fuxi-2000","version":"0.22.6","channels":"4020212223303800","moniker":"name","other":["amino_version=0.10.1","p2p_version=0.5.0","consensus_version=v1/0.2.2","rpc_version=0.7.0/3","tx_index=on","rpc_addr=tcp://0.0.0.0:26657"]},"sync_info":{"latest_block_hash":"CC9BBE0B38643DAF3D9B78D928E2ACA654E5A39C","latest_app_hash":"56B9228A97D5B85BFDBEE020E597D45D427ABC43","latest_block_height":"30048","latest_block_time":"2018-08-02T08:23:44.566550056Z","catching_up":false},"validator_info":{"address":"F638F7EA8A8E4DA559A346E1C404F83941749713","pub_key":{"type":"tendermint/PubKeyEd25519","value":"oI16LfBmnP8CefSGwIjAIO3QZ05xwB1+s4oPIQ3Yaag="},"voting_power":"10"}}
```
When you see `catching_up` is `false`, it means the node is synced with the rest of testnet, otherwise it means it's still syncing.
# Setup A Sentry Node

## Why do we need a sentry node?

A validator node is under the risk of distributed denial-of-service (DDoS)attack. It occurs when an attacker tries to disrupt normal traffic of a node. In this way, this node will be isolated from other nodes in the network. One way to mitigate this risk is for validators to carefully structure their network topology in a so-called sentry node architecture.

## What is a sentry node?

In IRISnet, a sentry node is just a  full node. The validator node will only connect to its sentry node. In this way, the sentry nodes will be protect the validator node from DDoS attack. 

## How to setup a sentry node?


### Sentry Node

On the sentry node's side, you need to get fully initialized first. 

Then, you should edit its `config.toml` file, and change `private_peers_id` field：

```
private_peer_ids =<validator_node_id>
```

`validator node id` is the `node-id` of validator node. 

Then you could start your sentry node,

```
iris init --moniker=<sentry name> --home=<sentry home>
```

```
{
  "chain_id": "test-chain",
  "node_id": "937efdf8526e3d9e8b5e887fa953ff1645cc096d",
  "app_message": {
    "secret": "issue envelope dose rail busy glass treat crop royal resemble city deer hungry govern cable angle cousin during mountain december spare stick unveil great"
  }
}
```

If you have multiple sentry node, you could make them as `persistent-peers` to each other. 

### Validator Node

On the validator node's side, you also need to get fully initialized first, and make sure you have the `priv_validator.json` file backuped. 

Then, you should edit its `config.toml` file,

```
persistent_peers = <sentry node id@sentry listen address:port">
persistent_peers = "bccb80864e540af719fa1be6ae913ce50a74feaf@ip:26656" 
```

If you want to put multiple sentry info, you need to separate the information with `,`. And the port is 26656 by default.

Set 
```
pex=false
``` 
In this way, the validator node will disable its peer reactor, so it will not respond to any peer exchange request other than its sentry nodes. 

Then you could start your validator node,

```
iris start --home=<validator node home>
```

It's also recommended to enable the firewall of validator node.  

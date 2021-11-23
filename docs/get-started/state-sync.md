---
order: 7
---

# State Sync

If you want to quickly start the node and join IRIS Hub without historical data, you can consider using the `state_sync` feature. **Note that the data directory must be empty when starting the node.**

## Procedure

1. Initialize the mainnet node by referring to [Join The Mainnet](./mainnet.md).
2. Check the block height of the current snapshot, and select the latest height.

```bash
curl http://34.82.96.8:26658/
```

3. Modify the `config.toml`.

```toml
[statesync]
enable = true # whether enable stat_sync; set true
rpc_servers = "34.82.96.8:26657,34.77.68.145:26657" # RPC server address which the node connects to
trust_height = # Set to the block height of the latest snapshot
trust_hash = "" #Set to the hash corresponding to the latest snapshot block height (trust height), which can be checked via https://irishub.iobscan.io/#/block/<trust_height>.
trust_period = "168h0m0s"
discovery_time = "15s"
temp_dir = ""
```

4. Start the node.

```bash
iris start
```

## Others

1. If any problem occurs during chain starting, you can execute `iris unsafe-reset-all` to reset the node and repeat the steps above.
2. If you can't find solutions to the current issue, please contact us via [IRISnet Discord channel](https://discord.com/invite/bmhu9F9xbX) for help.

# Status

Query node status

**Flags:**

| Name, shorthand | Default               | Description                                                | Required |
| --------------- | --------------------- | ---------------------------------------------------------- | -------- |
| --node, -n      | tcp://localhost:26657 | \<host>:\<port> to tendermint rpc interface for this chain |          |

## Query node status

```bash
iris status
```

Example Output:

```json
{
  "node_info": {
    "protocol_version": {
      "p2p": "7",
      "block": "10",
      "app": "0"
    },
    "id": "5e4f9d442a612566243f4b209752c73333a1511b",
    "listen_addr": "tcp://0.0.0.0:26656",
    "network": "test",
    "version": "0.33.6",
    "channels": "4020212223303800",
    "moniker": "node0",
    "other": {
      "tx_index": "on",
      "rpc_address": "tcp://0.0.0.0:26657"
    }
  },
  "sync_info": {
    "latest_block_hash": "5F22EDBAA0648FBAB6801ED5F553036E5B045606168A35839B20D55B9F6E06F3",
    "latest_app_hash": "537897D31AFC0E8B0357E86EE18E112FC33038CE471E4C715C6C414A2ADB6761",
    "latest_block_height": "3124",
    "latest_block_time": "2020-08-27T03:00:51.338539Z",
    "earliest_block_hash": "44452FEB0B8C1603EC62FC0336077B5159EA1574A0D01E88016225B8D3E38670",
    "earliest_app_hash": "",
    "earliest_block_height": "1",
    "earliest_block_time": "2020-08-26T06:43:07.065305Z",
    "catching_up": false
  },
  "validator_info": {
    "address": "AEBCB30D923D8410B4DE616152ACCF4DB2376351",
    "pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "LvMBCBylUBMyFWCZxh6AYymUvygUYNy8MKtZZeGe9xs="
    },
    "voting_power": "100"
  }
}
```

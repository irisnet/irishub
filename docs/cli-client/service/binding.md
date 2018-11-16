# iriscli service binding

## Description

Query service binding

## Usage

```
iriscli service binding [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                         | Required |
| --------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --bind-chain-id |                            | [string] the ID of the blockchain bond of the service               | Yes      |
| --def-chain-id  |                            | [string] the ID of the blockchain defined of the service            | Yes      |
| --provider      |                            | [string] bech32 encoded account created the service binding         | Yes      |
| --service-name  |                            | [string] service name                                               | Yes      |
| --help, -h      |                            | help for binding                                                    |          |
| --chain-id      |                            | [string] Chain ID of tendermint node                                |          |
| --height        | most recent provable block | [int] block height to query                                         |          |
| --indent        |                            | Add indent to JSON response                                         |          |
| --ledger        |                            | Use a connected Ledger device                                       |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node    | true                       | Don't verify proofs for responses                                   |          |


## Examples

### Query a service binding

```shell
iriscli service binding --def-chain-id=test --service-name=test-service --bind-chain-id=test --provider=faa1ydhmma8l4m9dygsh7l08fgrwka6yczs0gkfnvd
```

After that, you will get detail info for the service binding.

```json
{
  "type": "iris-hub/service/SvcBinding",
  "value": {
    "def_name": "test-service",
    "def_chain_id": "test",
    "bind_chain_id": "test",
    "provider": "faa1ydhmma8l4m9dygsh7l08fgrwka6yczs0gkfnvd",
    "binding_type": "Local",
    "deposit": [
      {
        "denom": "iris-atto",
        "amount": "1000000000000000000000"
      }
    ],
    "price": [
      {
        "denom": "iris-atto",
        "amount": "1000000000000000000"
      }
    ],
    "level": {
      "avg_rsp_time": "10000",
      "usable_time": "100"
    },
    "available": true,
    "disable_height": "0"
  }
}
```
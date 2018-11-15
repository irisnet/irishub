# iriscli stake delegations

## Description

Query all delegations made from one delegator

## Usage

```
iriscli stake delegations [delegator-addr] [flags]
```

## Flags

| Name, shorthand       | Default                    | Description                                                          | Required |
| --------------------- | -------------------------- | -------------------------------------------------------------------- | -------- |
| --chain-id            |                            | [string] Chain ID of tendermint node                                 |          |
| --height              | most recent provable block | block height to query                                                |          |
| --help, -h            |                            | help for delegations                                                 |          |
| --indent              |                            | Add indent to JSON response                                          |          |
| --ledger              |                            | Use a connected Ledger device                                        |          |
| --node                | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain  |          |
| --trust-node          | true                       | Don't verify proofs for responses                                    |          |

## Examples

### Query a validator

```shell
iriscli stake delegations DelegatorAddress
```

After that, you will get all detailed info of delegations from the specified delegator address.

```json
[
  {
    "delegator_addr": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
    "validator_addr": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
    "shares": "0.2000000000",
    "height": "290"
  }
]
```

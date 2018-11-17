# iriscli stake unbonding-delegations

## Description

Query all unbonding-delegations records for one delegator

## Usage

```
iriscli stake unbonding-delegations [delegator-addr] [flags]
```

## Flags

| Name, shorthand     | Default                    | Description                                                         | Required |
| ------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id          |                            | [string] Chain ID of tendermint node                                |          |
| --height            | most recent provable block | block height to query                                               |          |
| --help, -h          |                            | help for unbonding-delegations                                                  |          |
| --indent            |                            | Add indent to JSON response                                         |          |
| --ledger            |                            | Use a connected Ledger device                                       |          |
| --node              | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node        | true                       | Don't verify proofs for responses                                   |          |

## Examples

### Query an unbonding-delegation

```shell
iriscli stake unbonding-delegations DelegatorAddress
```

After that, you will get unbonding delegation's detailed info from specified delegator.

```json
[
  {
    "delegator_addr": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
    "validator_addr": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
    "creation_height": "1310",
    "min_time": "2018-11-15T06:24:22.754703377Z",
    "initial_balance": "0.02iris",
    "balance": "0.02iris"
  }
]
```

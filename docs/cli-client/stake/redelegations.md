# iriscli stake redelegations

## Description

Query all redelegations records for one delegator

## Usage

```
iriscli stake redelegations [delegator-addr] [flags]
```

## Flags

| Name, shorthand            | Default                    | Description                                                         | Required |
| -------------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id                 |                            | [string] Chain ID of tendermint node                                |          |
| --height                   | most recent provable block | block height to query                                               |          |
| --help, -h                 |                            | help for validator                                                  |          |
| --indent                   |                            | Add indent to JSON response                                         |          |
| --ledger                   |                            | Use a connected Ledger device                                       |          |
| --node                     | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node               | true                       | Don't verify proofs for responses                                   |          |

## Examples

### Query all redelegations records

```shell
iriscli stake redelegations DelegatorAddress
```

After that, you will get all redelegations records' info for specified delegator

```txt
TODO
```

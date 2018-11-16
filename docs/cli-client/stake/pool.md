# iriscli stake pool

## Description

Query the current staking pool values

## Usage

```
iriscli stake pool [flags]
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

### Query the current staking pool values

```shell
iriscli stake pool
```

After that, you will get the current staking pool values.

```txt
Pool
Loose Tokens: 49.8289125612
Bonded Tokens: 100.1800000000
Token Supply: 150.0089125612
Bonded Ratio: 0.6678269863
```

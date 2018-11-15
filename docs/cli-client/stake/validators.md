# iriscli stake validators

## Description

Query for all validators

## Usage

```
iriscli stake validators [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                         | Required |
| --------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                |          |
| --height        | most recent provable block | block height to query                                               |          |
| --help, -h      |                            | help for validator                                                  |          |
| --indent        |                            | Add indent to JSON response                                         |          |
| --ledger        |                            | Use a connected Ledger device                                       |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node    | true                       | Don't verify proofs for responses                                   |          |

## Examples

### Query a validator

```shell
iriscli stake validators
```

After that, you will get all validators' info.

```txt
Validator
Operator Address: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Validator Consensus Pubkey: fvp1zcjduepq47906n2zvq597vwyqdc0h35ve0fcl64hwqs9xw5fg67zj4g658aqyuhepj
Jailed: false
Status: Bonded
Tokens: 100.0000000000
Delegator Shares: 100.0000000000
Description: {node0   }
Bond Height: 0
Unbonding Height: 0
Minimum Unbonding Time: 1970-01-01 00:00:00 +0000 UTC
Commission: {{0.0000000000 0.0000000000 0.0000000000 0001-01-01 00:00:00 +0000 UTC}}
```

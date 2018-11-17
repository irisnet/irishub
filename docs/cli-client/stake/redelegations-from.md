# iriscli stake redelegations-from

## Description

Query all outgoing redelegatations from a validator

## Usage

```
iriscli stake redelegations-from [operator-addr] [flags]
```

## Flags

| Name, shorthand     | Default                    | Description                                                         | Required |
| ------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id          |                            | [string] Chain ID of tendermint node                                |          |
| --height            | most recent provable block | block height to query                                               |          |
| --help, -h          |                            | help for redelegations-from                                         |          |
| --indent            |                            | Add indent to JSON response                                         |          |
| --ledger            |                            | Use a connected Ledger device                                       |          |
| --node              | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node        | true                       | Don't verify proofs for responses                                   |          |

## Examples

### Query all outgoing redelegatations

```shell
iriscli stake redelegations-from ValidatorAddress
```

After that, you will get all outgoing redelegatations' from specified validator

```json
[
  {
    "delegator_addr": "faa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9sytjfk",
    "validator_src_addr": "fva1dayujdfnxjggd5ydlvvgkerp2supknthajpch2",
    "validator_dst_addr": "fva1h27xdw6t9l5jgvun76qdu45kgrx9lqede8hpcd",
    "creation_height": "1130",
    "min_time": "2018-11-16T07:22:48.740311064Z",
    "initial_balance": {
      "denom": "iris-atto",
      "amount": "100000000000000000"
    },
    "balance": {
      "denom": "iris-atto",
      "amount": "100000000000000000"
    },
    "shares_src": "100000000000000000.0000000000",
    "shares_dst": "100000000000000000.0000000000"
  }
]
```

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
| --help, -h          |                            | help for validator                                                  |          |
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

```txt
TODO
```

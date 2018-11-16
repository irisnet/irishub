# iriscli stake delegation

## Description

Query a delegation based on address and validator address

## Usage

```
iriscli stake delegation [flags]
```

## Flags

| Name, shorthand       | Default                    | Description                                                          | Required |
| --------------------- | -------------------------- | -------------------------------------------------------------------- | -------- |
| --address-delegator   |                            | [string] Bech address of the delegator                               | Yes      |
| --address-validator   |                            | [string] Bech address of the validator                               | Yes      |
| --chain-id            |                            | [string] Chain ID of tendermint node                                 |          |
| --height              | most recent provable block | block height to query                                                |          |
| --help, -h            |                            | help for delegation                                                   |          |
| --indent              |                            | Add indent to JSON response                                          |          |
| --ledger              |                            | Use a connected Ledger device                                        |          |
| --node                | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain  |          |
| --trust-node          | true                       | Don't verify proofs for responses                                    |          |

## Examples

### Query a validator

```shell
iriscli stake delegation --address-validator=ValidatorAddress --address-delegator=DelegatorAddress

```

After that, you will get detailed info of the delegation between specified validator and delegator.

```txt
Delegation
Delegator: faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx
Validator: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Shares: 0.2000000000Height: 290
```

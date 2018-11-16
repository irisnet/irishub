# iriscli stake unbonding-delegation

## Description

Query an unbonding-delegation record based on delegator and validator address

## Usage

```
iriscli stake unbonding-delegation [flags]
```

## Flags

| Name, shorthand     | Default                    | Description                                                         | Required |
| ------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --address-delegator |                            | [string] Bech address of the delegator                              | Yes      |
| --address-validator |                            | [string] Bech address of the validator                              | Yes      |
| --chain-id          |                            | [string] Chain ID of tendermint node                                |          |
| --height            | most recent provable block | block height to query                                               |          |
| --help, -h          |                            | help for unbonding-delegation                                       |          |
| --indent            |                            | Add indent to JSON response                                         |          |
| --ledger            |                            | Use a connected Ledger device                                       |          |
| --node              | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node        | true                       | Don't verify proofs for responses                                   |          |


## Examples

### Query an unbonding-delegation

```shell
iriscli stake unbonding-delegation --address-delegator=DelegatorAddress --address-validator=ValidatorAddress
```

After that, you will get unbonding delegation's detailed info between specified validator and delegator.

```txt
Unbonding Delegation
Delegator: faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx
Validator: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Creation height: 1310
Min time to unbond (unix): 2018-11-15 06:24:22.754703377 +0000 UTC
Expected balance: 0.02iris
```

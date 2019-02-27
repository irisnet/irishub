# iriscli stake unbonding-delegation

## Description

Query an unbonding-delegation record based on delegator and validator address

## Usage

```
iriscli stake unbonding-delegation [flags]
```
Print help messages:
```
iriscli stake unbonding-delegation --help
```

## Unique Flags

| Name, shorthand     | Default                    | Description                                                         | Required |
| ------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --address-delegator |                            | [string] Bech address of the delegator                              | Yes      |
| --address-validator |                            | [string] Bech address of the validator                              | Yes      |


## Examples

Query an unbonding-delegation
```
iriscli stake unbonding-delegation --address-delegator=DelegatorAddress --address-validator=ValidatorAddress
```

After that, you will get unbonding delegation's detailed info between specified validator and delegator.

```txt
Unbonding Delegation
Delegator: iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
Validator: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Creation height: 1310
Min time to unbond (unix): 2018-11-15 06:24:22.754703377 +0000 UTC
Expected balance: 0.02iris
```

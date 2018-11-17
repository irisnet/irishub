# iriscli stake delegation

## Description

Query a delegation based on address and validator address

## Usage

```
iriscli stake delegation [flags]
```
Print help messages:
```
iriscli stake delegation --help
```

## Unique Flags

| Name, shorthand       | Default                    | Description                                                          | Required |
| --------------------- | -------------------------- | -------------------------------------------------------------------- | -------- |
| --address-delegator   |                            | [string] Bech address of the delegator                               | Yes      |
| --address-validator   |                            | [string] Bech address of the validator                               | Yes      |

## Examples

Query a validator
```
iriscli stake delegation --address-validator=ValidatorAddress --address-delegator=DelegatorAddress

```

After that, you will get detailed info of the delegation between specified validator and delegator.

```txt
Delegation
Delegator: faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx
Validator: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Shares: 0.2000000000Height: 290
```

# iriscli stake unbonding-delegations

## Description

Query all unbonding-delegations records for one delegator

## Usage

```
iriscli stake unbonding-delegations [delegator-address] [flags]
```
Print help messages:
```
iriscli stake unbonding-delegations --help
```

## Examples

Query an unbonding-delegation
```
iriscli stake unbonding-delegations [delegator-address]
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

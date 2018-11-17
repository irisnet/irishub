# iriscli stake unbonding-delegations-from

## Description

Query all unbonding delegatations from a validator

## Usage

```
iriscli stake unbonding-delegations-from [validator-address] [flags]
```
Print help messages:
```
iriscli stake unbonding-delegations-from --help
```

## Examples

Query all unbonding delegatations from a validator
```
iriscli stake unbonding-delegations [validator-address]
```

After that, you will get unbonding delegation's detailed info from specified validator.

```json
[
  {
    "delegator_addr": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
    "validator_addr": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
    "creation_height": "1310",
    "min_time": "2018-11-15T06:24:22.754703377Z",
    "initial_balance": {
      "denom": "iris-atto",
      "amount": "20000000000000000"
    },
    "balance": {
      "denom": "iris-atto",
      "amount": "20000000000000000"
    }
  }
]
```

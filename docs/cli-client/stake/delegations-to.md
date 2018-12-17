# iriscli stake delegations-to

## Description

Query all delegations to one validator

## Usage

```
iriscli stake delegations-to [validator-address] [flags]
```
Print help messages:
```
iriscli stake delegations-to --help
```

## Examples

Query all delegations to one validator
```
iriscli stake delegations-to [validator-address]
```

After that, you will get all detailed info of delegations from the specified delegator address.

```json
[
  {
    "delegator_addr": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
    "validator_addr": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
    "shares": "0.2000000000",
    "height": "290"
  }
]
```

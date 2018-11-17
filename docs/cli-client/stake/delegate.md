# iriscli stake delegate

## Introduction

Delegate tokens to a validator

## Usage

```
iriscli stake delegate [flags]
```

Print help messages:
```
iriscli stake delegate --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-delegator | string | true     | ""       | Bech address of the delegator |
| --amount            | string | true     | ""       | Amount of coins to bond |

## Examples

```
iriscli stake delegate --chain-id=<chain-id> --from=<key name> --fee=0.004iris --amount=100iris --address-validator=<ValidatorAddress>
```

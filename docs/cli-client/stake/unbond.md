# iriscli stake unbond

## Introduction

Unbond shares from a validator

## Usage

```
iriscli stake unbond [flags]
```

Print all help messages:

```shell
iriscli stake unbond --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator | string | true     | ""       | Bech address of the validator |
| --shares-amount     | float  | false    | 0.0      | Amount of source-shares to either unbond or redelegate as a positive integer or decimal |
| --shares-percent    | float  | false    | 0.0      | Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 |

Users must specify the unbond amount. There two options can do this: `--shares-amount` or `--shares-percent`. Keep in mind, don't specify them both.

## Examples

```shell
iriscli stake unbond --address-validator=<ValidatorAddress> --shares-percent=0.1 --from=<key name> --chain-id=<chain-id> --fee=0.004iris
```

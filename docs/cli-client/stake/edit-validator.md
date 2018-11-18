# iriscli stake edit-validator

## Introduction

Edit existing validator, such as commission rate, name and other description message.

## Usage

```
iriscli stake edit-validator [flags]
```
Print help messages:
```
iriscli stake edit-validator --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --commission-rate   | string | float    | 0.0      | Commission rate percentage |
| --moniker           | string | false    | ""       | Validator name |
| --identity          | string | false    | ""       | Optional identity signature (ex. UPort or Keybase) |
| --website           | string | false    | ""       | Optional website  |
| --details           | string | false    | ""       | Optional details |


## Examples

```
iriscli stake edit-validator --from=<key name> --chain-id=<chain-id> --fee=0.004iris --commission-rate=0.15
```

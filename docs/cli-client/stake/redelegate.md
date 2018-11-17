# iriscli stake redelegate

## Introduction

Transfer delegation from one validator to another one.

## Usage

```
iriscli stake redelegate [flags]
```

Print all help messages:

```shell
iriscli stake redelegate --help
```

## Unique Flags

| Name, shorthand            | type   | Required | Default  | Description                                                         |
| -------------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator-dest   | string | true     | ""       | Bech address of the destination validator |
| --address-validator-source | string | true     | ""       | Bech address of the source validator |
| --shares-amount            | float  | false    | 0.0      | Amount of source-shares to either unbond or redelegate as a positive integer or decimal |
| --shares-percent           | float  | false    | 0.0      | Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 |

Users must specify the redeleagtion token amount. There two options can do this: `--shares-amount` or `--shares-percent`. Keep in mind, don't specify both of them.

## Examples

```shell
iriscli stake redelegate --chain-id=<chain-id> --from=<key name> --fee=0.004iris --address-validator-source=<SourceValidatorAddress> --address-validator-dest=<DestinationValidatorAddress> --shares-percent=0.1
```

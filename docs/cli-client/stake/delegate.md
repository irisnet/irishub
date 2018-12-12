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
> Note: bech32 prefix for validator address is `fva` 


For example, if you want to delegate 10iris on fuxi-6000:
```$xslt
iriscli stake delegate --chain-id=fuxi-6000 --from=abc --fee=0.004iris --amount=10iris --address-validator=fva12zgt9hc5r5mnxegam9evjspgwhkgn4wzjxkvqy
```
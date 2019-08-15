# iriscli bank burn

## Description

This command is used to burn tokens from your own address

## Usage

```bash
iriscli bank burn --from=<key-name> --amount=<amount-to-burn> --fee=<native-fee> --chain-id=<chain-id>
```

## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --amount         | String | True     |                       | Amount of coins to burn, e.g. 10iris                         |

## Examples

### Burn Token

```bash
 iriscli bank burn --from=<key-name> --amount=10iris --chain-id=irishub --fee=0.3iris
```

# iriscli bank set-memo-regexp

## Description

This command is used to set memo regexp for your own address

## Usage

```bash
iriscli bank set-memo-regexp --regexp=<regular-expression> --from=<key-name> --fee=<native-fee> --chain-id=<chain-id>
```

## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --regexp         | String | True     |                       | Regular expression, maximum length 50, e.g. ^[A-Za-z0-9]+$   |

## Examples

### Send tokens to another address

```bash
iriscli bank set-memo-regexp --regexp=^[A-Za-z0-9]+$ --from=<key-name> --amount=10iris --fee=0.3iris --chain-id=irishub
```

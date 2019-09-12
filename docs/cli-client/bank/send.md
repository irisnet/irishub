# iriscli bank send

## Description

Sending tokens to another address, this command includes `generate/sign/broadcast` transactions.

## Usage

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=<amount> --fee=<native-fee> --chain-id=<chain-id>
```

## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --amount         | string | true     |                       | Amount of coins to send, for instance: 10iris                |
| --to             | string |          |                       | Bech32 encoding address to receive coins                     |

## Examples

### Send tokens to another address

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10iris --fee=0.3iris --chain-id=irishub
```

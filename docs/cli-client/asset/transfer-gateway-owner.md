# iriscli asset transfer-gateway-owner

## Introduction

Transfer the ownership of a gateway to a new owner.

## Usage

```bash
iriscli asset transfer-gateway-owner [flags]
```

## Unique Flags

| Name, shorthand     | type   | Required | Default   | Description                                                       |
| --------------------| -----  | -------- | --------  |-------------------------------------------------------- |
| --moniker           | string  | true     | ""       | the unique name of the gateway to be transferred       |
| --to                | Address | true     |          | the new owner to which the gateway will be transferred |

## Examples

```bash
iriscli asset transfer-gateway-owner --moniker=cats --to=<new-owner-address> --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

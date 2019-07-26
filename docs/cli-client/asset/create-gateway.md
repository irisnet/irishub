# iriscli asset create-gateway

## Introduction

Create a gateway which is used to map external assets

## Usage

```bash
iriscli asset create-gateway [flags]
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | true     | ""       | The unique name with a size between 3 and 8, beginning with a letter followed by alphanumeric characters|
| --identity          | string | false    | ""       | Optional identity signature with a maximum length of 128 (ex. UPort or Keybase)|
| --details           | string | false    | ""       | Optional details with a maximum length of 280|
| --website           | string | false    | ""       | Optional website with a maximum length of 128|

## Examples

```bash
iriscli asset create-gateway --moniker=cats --identity=<pgp-id> --details="Cat Tokens" --website="www.example.com" --from=<key-name> --chain-id=irishub --fee=0.3iris
```

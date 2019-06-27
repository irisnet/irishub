# iriscli asset query-fee

## Introduction

Query the asset-related fees, including gateway creation and token issuance and minting

## Usage

```
iriscli asset query-fee [flags]
```

Print help messages:
```
iriscli asset query-fee --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --subject           | string | true     | ""       | the fee type, only "gateway" and "token" allowed       |
| --moniker           | string | false    | ""       | the gateway name, required if the subject is "gateway" |
| --id                | string | false    | ""       | the token id, required if the subject is "token"       |


## Examples

```
iriscli asset query-fee --subject gateway --moniker tgw
```

Output:
```txt
Fee: 600000iris
```

```json
{
  "exist": false,
  "fee": {
    "denom": "iris",
    "amount": "600000"
  }
}
```
# iriscli asset query-fee

## Introduction

Query the asset related fees, including gateway creation and token issuance and minting

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
| --gateway           | string | false    | ""       | the gateway moniker, required for querying gateway fee |
| --token             | string | false    | ""       | the token id, required for querying token fees         |


## Examples

```
iriscli asset query-fee --gateway=tgw
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

```
iriscli asset query-fee --token=i.sym
```

Output:
```txt
Fees:
  IssueFee: 300000iris
  MintFee:  30000iris
```

```json
{
  "Exist": false,
  "issue_fee": {
    "denom": "iris",
    "amount": "300000"
  },
  "mint_fee": {
    "denom": "iris",
    "amount": "30000"
  }
}
```
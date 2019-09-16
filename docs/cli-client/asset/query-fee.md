# iriscli asset query-fee

## Introduction

Query the asset related fees, including gateway creation and token issuance and minting

## Usage

```bash
iriscli asset query-fee [flags]
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                            |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------ |
| --gateway           | string |          |          | The gateway moniker, required for querying gateway fee |
| --token             | string |          |          | The token id, required for querying token fees         |

## Examples

### Query fee of creating a gateway

```bash
iriscli asset query-fee --gateway=cats
```

### Query fee of issuing and minting a native token

```bash
iriscli asset query-fee --token=kitty
```

### Query fee of issuing and minting a gateway token

```bash
iriscli asset query-fee --token=cats.kitty
```

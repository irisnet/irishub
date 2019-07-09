# iriscli asset query-gateway

## Introduction

Query a gateway by moniker

## Usage

```bash
iriscli asset query-gateway [flags]
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | true     | ""       | the unique name with a size between 3 and 8 letters|

## Examples

```bash
iriscli asset query-gateway --moniker cats
```

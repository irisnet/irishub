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
| --moniker           | string | true     | ""       | The unique name with a size between 3 and 8, beginning with a letter followed by alphanumeric characters |

## Examples

```bash
iriscli asset query-gateway --moniker cats
```

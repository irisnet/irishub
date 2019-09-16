# iriscli rand query-rand

## Introduction

Query the generated random number by the request id

## Usage

```bash
iriscli rand query-rand [flags]
```

## Unique Flags

| Name, shorthand     | type   | Required | Default | Description                               |
| --------------------| ------ | -------- | ------- | ----------------------------------------- |
| --request-id        | string | true     |         | The request id returned by the request tx |

## Examples

```bash
iriscli rand query-rand --request-id=035a8d4cf64fcd428b5c77b1ca85bfed172d3787be9bdf0887bbe8bbeec3932c
```

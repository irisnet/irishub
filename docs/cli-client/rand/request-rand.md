# iriscli rand request-rand

## Introduction

Request a random number

## Usage

```bash
iriscli rand request-rand [flags]
```

## Unique Flags

| Name, shorthand     | type   | Required | Default | Description                                                                  |
| --------------------| ------ | -------- | ------- | ---------------------------------------------------------------------------- |
| --block-interval    | uint64 |          | 10      | The block interval after which the requested random number will be generated |

## Examples

```bash
iriscli rand request-rand --block-interval=100 --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```

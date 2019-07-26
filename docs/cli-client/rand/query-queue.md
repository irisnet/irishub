# iriscli rand query-queue

## Introduction

Query the pending random number requests with an optional height

## Usage

```bash
iriscli rand query-queue [flags]
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --queue-height      | int64  | false     |  0      | the height at which the pending requests will be retrieved |

## Examples

```bash
iriscli rand query-queue --queue-height=100000
```

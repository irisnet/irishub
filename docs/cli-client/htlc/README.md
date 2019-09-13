# iriscli htlc

## Description

Htlc allows you to manage local htlc for atomic swap with other chains

## Usage

```
 iriscli htlc [command]
```


## Available Commands

| Name                        | Description                 |
| --------------------------- | --------------------------- |
| [create](create.md)         | Create a HTLC               |
| [claim](claim.md)           | Claim an opened HTLC        |
| [refund](refund.md)         | Refund from an expired HTLC |
| [query-htlc](query-htlc.md) | Query details of a htlc     |

## Global Flags

| Name, shorthand | Default        | Description                                 | Required | Type   |
| --------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding  | hex            | String binary encoding (hex\b64\btc )       |          | string |
| --home          | /root/.iriscli | Directory for config and data               |          | string |
| -o, --output    | text           | Output format (text\json)                   |          | string |
| --trace         |                | Print out full stack trace on errors        |          |        |

## Flags

| Name, shorthand | Default | Description   | Required |
| --------------- | ------- | ------------- | -------- |
| -h, --help      |         | Help for htlc |          |

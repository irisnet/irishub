# iriscli bank broadcast

## Description

This command is used for broadcasting a signed transaction to the network after generating a signed transaction offline with [sign](./sign.md).

## Usage:

```
iriscli bank broadcast <tx-file> <flags> 
```

## Flags

| Name,shorthand | Type   | Required | Default               | Description                                                  |
| -------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help     |        | False    |                       | help for account                                             |
| --chain-id     | String | False    |                       | Chain ID of tendermint node                                  |
| --height       | Int    | False    |                       | Block height to query, omit to get most recent provable block |
| --ledger       | String | False    |                       | Use a connected Ledger device                                |
| --node         | String | False    | tcp://localhost:26657 | <host>:<port> to tendermint rpc interface for this chain     |
| --trust-node   | String | False    | True                  | Don't verify proofs for responses                            |

## Global Flags

| Name,shorthand        | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding string | hex            | String   Binary encoding (hex\b64\btc ) | False    | String |
| --home string         | /root/.iriscli | Directory for config and data               | False    | String |
| -o, --output string   | text           | Output format (text\json)                 | False    | String |
| --trace               |                | Print out full stack trace on errors        | False    |        |

## Examples

### Broadcast your transaction

```
iriscli bank broadcast sign.json --chain-id=<chain-id>
```

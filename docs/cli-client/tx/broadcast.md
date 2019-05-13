# iriscli tx broadcast

## Description

This command is used for broadcasting a signed transaction to the network after generating a signed transaction offline with [sign](./sign.md).

## Usage:

```
iriscli tx broadcast <tx-file> <flags> 
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

## Examples

### Broadcast your transaction

```
iriscli tx broadcast sign.json --chain-id=<chain-id>
```

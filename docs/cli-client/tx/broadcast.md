# iriscli tx broadcast

## Description

This command is used for broadcasting a signed transaction to the network after generating a signed transaction offline with [sign](./sign.md).

## Usage:

```
iriscli tx broadcast <tx-file> <flags> 
```

## Flags

| Name, shorthand | Type   | Required | Default               | Description                                                   |
| --------------- | ------ | -------- | --------------------- | ------------------------------------------------------------- |
| -h, --help      |        |          |                       | Help for account                                              |
| --chain-id      | string |          |                       | Chain ID of tendermint node                                   |
| --height        | int    |          |                       | Block height to query, omit to get most recent provable block |
| --ledger        | string |          |                       | Use a connected Ledger device                                 |
| --node          | string |          | tcp://localhost:26657 | `<host>:<port>`to tendermint rpc interface for this chain     |
| --trust-node    | string |          | true                  | Don't verify proofs for responses                             |

## Examples

### Broadcast your transaction

```
iriscli tx broadcast sign.json --chain-id=<chain-id>
```

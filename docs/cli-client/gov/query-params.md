# iriscli gov query-params

## Description

Query parameter proposal's config

## Usage

```
iriscli gov query-params [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | help for submit-proposal                                                                                                                             |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --key           |                            | [string] key name of parameter                                                                                                                       | Yes      |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --module        |                            | [string] module name                                                                                                                                 |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

## Examples
 
### Query params

```shell

```

 After that, you're done with depositing iris tokens for an activing proposal, and remember to back up your proposal-id, it's the only way to retrieve your proposal.

```txt

```
# iriscli gov query-tally

## Description

Get the tally of a proposal vote
 
## Usage

```
iriscli gov query-tally [flags]
```

## Flags
| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | help for submit-proposal                                                                                                                             |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                            | [string] proposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

## Examples

### Query tally

```shell

```
 After that, you're done with depositing iris tokens for an activing proposal, and remember to back up your proposal-id, it's the only way to retrieve your proposal.

```txt

```
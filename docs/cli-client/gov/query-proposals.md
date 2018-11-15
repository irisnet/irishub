# iriscli gov query-proposals

## Description

Query proposals with optional filters

## Usage

```
iriscli gov query-proposals [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --depositer     |                            | [string] (optional) filter by proposals deposited on by depositer                                                                                    |          |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | help for submit-proposal                                                                                                                             |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --limit         |                            | [string] (optional) limit to latest [number] proposals. Defaults to all proposals                                                                    |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --status        |                            | [string] proposalID of proposal depositing on                                                                                                        |          |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |
| --voter         |                            | [string] (optional) filter by proposals voted on by voted                                                                                            |          |

## Examples

### Query proposals

```shell
iriscli gov query-proposals --chain-id=test
```

 After that, you're done with depositing iris tokens for an activing proposal, and remember to back up your proposal-id, it's the only way to retrieve your proposal.

```txt
  1 - test proposal
  2 - new proposal
```
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
| --height        |                            | [int] Block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | Help for query-tally                                                                                                                                 |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                            | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

## Examples

### Query tally

```shell
iriscli gov query-tally --chain-id=test --proposal-id=1
```

You could query the statistics of each voting option.

```txt
{
  "yes": "100.0000000000",
  "abstain": "0.0000000000",
  "no": "0.0000000000",
  "no_with_veto": "0.0000000000"
}
```
# iriscli gov query-tally

## Description

Get the tally of a proposal vote
 
## Usage

```
iriscli gov query-tally <flags>
```

Print help messages:

```
iriscli gov query-tally --help
```

## Flags
| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | ProposalID of proposal depositing on                                                                                                        | Yes      |

## Examples

### Query tally

You could query the statistics of each voting option.

```shell
iriscli gov query-tally --chain-id=<chain-id> --proposal-id=<proposal-id>
```

```txt
Tally Result:
  Yes:        0
  Abstain:    100.0000000000
  No:         0
  NoWithVeto: 200.0000000000
```

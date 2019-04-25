# iriscli gov query-tally

## Description

Get the tally of a proposal vote
 
## Usage

```
iriscli gov query-tally [flags]
```

Print help messages:

```
iriscli gov query-tally --help
```

## Flags
| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |

## Examples

### Query tally

```shell
iriscli gov query-tally --chain-id=<chain-id> --proposal-id=1
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

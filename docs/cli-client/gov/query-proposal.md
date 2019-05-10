# iriscli gov query-proposal

## Description

Query details of a single proposal

## Usage

```
iriscli gov query-proposal <flags>
```

Print help messages:

```
iriscli gov query-proposal --help
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | ProposalID of proposal depositing on                                                                                                        | Yes      |

## Examples

### Query proposal

```shell
iriscli gov query-proposal --chain-id=<chain-id> --proposal-id=<proposal-id>
```

You could query the details of a specific proposal.

```txt
Proposal 94:
  Title:              test proposal
  Type:               TxTaxUsage
  Status:             Rejected
  Submit Time:        2019-05-10 06:37:18.776274942 +0000 UTC
  Deposit End Time:   2019-05-10 06:37:28.776274942 +0000 UTC
  Total Deposit:      1100iris
  Voting Start Time:  2019-05-10 06:37:18.776274942 +0000 UTC
  Voting End Time:    2019-05-10 06:37:28.776274942 +0000 UTC
  Description:        a new text proposal
```

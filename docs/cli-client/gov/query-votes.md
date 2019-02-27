# iriscli gov query-votes

## Description

Query votes on a proposal

## Usage

```
iriscli gov query-votes [flags]
```

Print help messages:

```
iriscli gov query-votes --help
```
## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |

## Examples

### Query votes

```shell
iriscli gov query-votes --chain-id=<chain-id> --proposal-id=1
```

You could query the voting of all the voters by specifying the proposal.
 
```txt
[
  {
    "voter": "iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fpascegs",
    "proposal_id": "1",
    "option": "Yes"
  }
]
```

# iriscli gov query-votes

## Description

Query votes on a proposal

## Usage

```
iriscli gov query-votes <flags>
```

Print help messages:

```
iriscli gov query-votes --help
```
## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | ProposalID of proposal depositing on                                                                                                        | Yes      |

## Examples

### Query votes

You could query the voting of all the voters by specifying the proposal id.
 
```shell
iriscli gov query-votes --chain-id=<chain-id> --proposal-id=<proposal-id>
```
 
```txt
Votes for Proposal 99:
  iaa1gfcee5u5f54kfcnufv4ypcfyldw0vu0z5l4mh8: Abstain
  iaa15x7ph3pz5dvh92dzhjfcrglswu2r9uygly5vmu: NoWithVeto
```

# iriscli gov query-vote

## Description

Query vote

## Usage

```
iriscli gov query-vote <flags>
```

Print help messages:

```
iriscli gov query-vote --help
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | ProposalID of proposal depositing on                                                                                                        | Yes      |
| --voter         |                            | Bech32 voter address                                                                                                                        | Yes      |

## Examples

### Query vote

You could query the voting by specifying the proposal id and the voter.

```shell
iriscli gov query-vote --chain-id=<chain-id> --proposal-id=<proposal-id> --voter=<voter_address>
```

```txt
{
  "voter": "iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fpascegs",
  "proposal_id": "1",
  "option": "Yes"
}
```

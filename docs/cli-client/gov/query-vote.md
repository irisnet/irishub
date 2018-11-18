# iriscli gov query-vote

## Description

Query vote

## Usage

```
iriscli gov query-vote [flags]
```

Print help messages:

```
iriscli gov query-vote --help
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |
| --voter         |                            | [string] Bech32 voter address                                                                                                                        | Yes      |

## Examples

### Query vote

```shell
iriscli gov query-vote --chain-id=test --proposal-id=1 --voter=faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd
```

You could query the voting by specifying the proposal and the voter.

```txt
{
  "voter": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd",
  "proposal_id": "1",
  "option": "Yes"
}
```

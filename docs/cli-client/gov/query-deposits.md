# iriscli gov query-deposits

## Description

Query details of a deposits

## Usage

```
iriscli gov query-deposits [flags]
```

Print help messages:

```
iriscli gov query-deposits --help
```
## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | ProposalID of proposal depositing on                                                                                                        | Yes      |


## Examples

### Query deposits

```shell
iriscli gov query-deposits --chain-id=<chain-id> --proposal-id=<proposal-id>
```

You could query all the deposited tokens on a specific proposal, includes deposit details for each depositor.

```txt
Deposits for Proposal 92:
  iaa1c4kjt586r3t353ek9jtzwxum9x9fcgwent790r: 995iris
```

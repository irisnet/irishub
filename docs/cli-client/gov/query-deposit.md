# iriscli gov query-deposit

## Description

Query details of a deposit

## Usage

```
iriscli gov query-deposit <flags>
```

Print help messages:

```
iriscli gov query-deposit --help
```
## Flags

| Name, shorthand | Default               | Description                                                                                                                                          | Required |
| --------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                       | Chain ID of tendermint node                                                                                                                 | Yes      |
| --depositor     |                       | Bech32 depositor address                                                                                                                    | Yes      |
 
## Examples

### Query deposit

```shell
iriscli gov query-deposit --chain-id=<chain-id> --proposal-id=<proposal-id> --depositor=<depositor_address>
```

You could query the deposited tokens on a specific proposal by `proposal-id` and `depositor`.

```txt
Deposit by iaa1c4kjt586r3t353ek9jtzwxum9x9fcgwent790r on Proposal 90 is for the amount 995iris
```

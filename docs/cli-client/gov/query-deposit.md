# iriscli gov query-deposit

## Description

Query details of a deposit

## Usage

```
iriscli gov query-deposit [flags]
```

Print help messages:

```
iriscli gov query-deposit --help
```
## Flags

| Name, shorthand | Default               | Description                                                                                                                                          | Required |
| --------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                       | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --depositor     |                       | [string] Bech32 depositor address                                                                                                                    | Yes      |
 
## Examples

### Query deposit

```shell
iriscli gov query-deposit --chain-id=<chain-id> --proposal-id=1 --depositor=faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07
```

You could query the deposited tokens on a specific proposal.

```txt
{
  "depositor": "faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07",
  "proposal_id": "1",
  "amount": [
    {
      "denom": "iris-atto",
      "amount": "30000000000000000000"
    }
  ]
}
```

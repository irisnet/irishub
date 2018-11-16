# iriscli gov query-deposit

## Description

Query details of a deposit

## Usage

```
iriscli gov query-deposit [flags]
```

## Flags

| Name, shorthand | Default               | Description                                                                                                                                          | Required |
| --------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                       | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --depositer     |                       | [string] Bech32 depositer address                                                                                                                    | Yes      |
| --height        |                       | [int] Block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                       | Help for query-deposit                                                                                                                               |          |
| --indent        |                       | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                       | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657 | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                       | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                  | Don't verify proofs for responses                                                                                                                    |          |
 
## Examples

### Query deposit

```shell
iriscli gov query-deposit --chain-id=test --proposal-id=1 --depositer=faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07
```

You could query the deposited tokens on a specific proposal.

```txt
{
  "depositer": "faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07",
  "proposal_id": "1",
  "amount": [
    {
      "denom": "iris-atto",
      "amount": "30000000000000000000"
    }
  ]
}
```
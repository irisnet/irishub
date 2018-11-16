# iriscli gov query-deposits

## Description

Query details of a deposits

## Usage

```
iriscli gov query-deposits [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --height        |                            | [int] Block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | Help for query-deposits                                                                                                                              |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                            | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

## Examples

### Query deposits

```shell
iriscli gov query-deposits --chain-id=test --proposal-id=1
```

You could query all the deposited tokens on a specific proposal, includes deposit details for each depositor.

```txt
[
  {
    "depositer": "faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07",
    "proposal_id": "1",
    "amount": [
      {
        "denom": "iris-atto",
        "amount": "35000000000000000000"
      }
    ]
  }
]
```
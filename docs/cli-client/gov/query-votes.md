# iriscli gov query-votes

## Description

Query votes on a proposal

## Usage

```
iriscli gov query-votes [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --height        |                            | [int] Block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | Help for query-votes                                                                                                                                 |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                            | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

## Examples

### Query votes

```shell
iriscli gov query-votes --chain-id=test --proposal-id=1
```

You could query the voting of all the voters by specifying the proposal.
 
```txt
[
  {
    "voter": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd",
    "proposal_id": "1",
    "option": "Yes"
  }
]
```
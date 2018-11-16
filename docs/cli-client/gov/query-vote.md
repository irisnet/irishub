# iriscli gov query-vote

## Description

Query vote

## Usage

```
iriscli gov query-vote [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | help for submit-proposal                                                                                                                             |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                            | [string] proposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |
| --voter         |                            | [string] bech32 voter address                                                                                                                        | Yes      |

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
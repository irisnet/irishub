# iriscli gov query-proposals

## Description

Query proposals with optional filters

## Usage

```
iriscli gov query-proposals [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                                                                                                          | Required |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --depositer     |                            | [string] (optional) Filter by proposals deposited on by depositer                                                                                    |          |
| --height        |                            | [int] Block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | Help for query-proposals                                                                                                                             |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --limit         |                            | [string] (optional) Limit to latest [number] proposals. Defaults to all proposals                                                                    |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --status        |                            | [string] (optional) filter proposals by proposal status                                                                                                        |          |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |
| --voter         |                            | [string] (optional) Filter by proposals voted on by voted                                                                                            |          |

## Examples

### Query proposals

```shell
iriscli gov query-proposals --chain-id=test
```

You could query all the proposals by default.

```txt
  1 - test proposal
  2 - new proposal
```

Also you can query proposal by filters, such as:

```shell
gov query-proposals --chain-id=test --depositer=faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd
```

Finally, here shows the proposal who's depositor address is faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd.

```txt
  2 - new proposal
```
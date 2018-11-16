# iriscli gov query-proposals

## 描述

通过可选项过滤查询满足条件的提议

## 使用方式

```
iriscli gov query-proposals [flags]
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --depositer     |                            | [string] (optional) filter by proposals deposited on by depositer                                                                                    |          |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | Help for query-proposals                                                                                                                             |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --limit         |                            | [string] (optional) limit to latest [number] proposals. Defaults to all proposals                                                                    |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --status        |                            | [string] proposalID of proposal depositing on                                                                                                        |          |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |
| --voter         |                            | [string] (optional) filter by proposals voted on by voted                                                                                            |          |

## 例子

### 查询提议

```shell
iriscli gov query-proposals --chain-id=test
```

默认查询所有的提议。

```txt
  1 - test proposal
  2 - new proposal
```

当然这里可以查询指定条件的提议。

```shell
gov query-proposals --chain-id=test --depositer=faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd
```

可以得到存款人是faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd地址的提议。
```txt
  2 - new proposal
```

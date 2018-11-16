# iriscli gov query-tally

## 描述

查询指定提议的投票统计
 
## 使用方式

```
iriscli gov query-tally [flags]
```

## 标志
| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | help for submit-proposal                                                                                                                             |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                            | [string] proposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

## 例子

### 查询投票统计

```shell
iriscli gov query-tally --chain-id=test --proposal-id=1
```

可以查询指定提议每个投票选项的投票统计。

```txt
{
  "yes": "100.0000000000",
  "abstain": "0.0000000000",
  "no": "200.0000000000",
  "no_with_veto": "0.0000000000"
}
```
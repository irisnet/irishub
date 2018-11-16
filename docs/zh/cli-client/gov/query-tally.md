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
| --chain-id      |                            | [string] tendermint节点的链ID                                                                                                                 | Yes      |
| --height        |                            | [int] 查询的区块高度                                                                                  |          |
| --help, -h      |                            | 查询命令帮助                                                                                                                                 |          |
| --indent        |                            | 在JSON响应中添加缩进                                                                                                                          |          |
| --ledger        |                            | 使用连接的硬件记账设备                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                                                                                  |          |
| --proposal-id   |                            | [string] 提议ID                                                                                                        | Yes      |
| --trust-node    | true                       | 关闭响应结果校验                                                                                                                    |          |

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
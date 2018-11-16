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
| --chain-id      |                            | [string] tendermint节点的链ID                                                                                                                 | Yes      |
| --depositer     |                            | [string] （可选）按存款人过滤                                                                                    |          |
| --height        |                            | [int] 查询的区块高度                                                                                  |          |
| --help, -h      |                            | 查询命令帮助                                                                                                                             |          |
| --indent        |                            | 在JSON响应中添加缩进                                                                                                                          |          |
| --ledger        |                            | 使用连接的硬件记账设备                                                                                                                        |          |
| --limit         |                            | [string] （可选）限制最新[数量]提议。 默认为所有提议                                                                    |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                                                                                  |          |
| --status        |                            | [string] （可选）按提议状态过滤提议                                                                                                        |          |
| --trust-node    | true                       | 关闭响应结果校验                                                                                                                    |          |
| --voter         |                            | [string] （可选）按投票人过滤                                                                                            |          |

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

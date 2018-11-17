# iriscli gov query-votes

## 描述

查询指定提议的投票情况

## 使用方式

```
iriscli gov query-votes [flags]
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

### Query votes

```shell
iriscli gov query-votes --chain-id=test --proposal-id=1
```

通过指定的提议查询该提议所有投票者的投票详情。
 
```txt
[
  {
    "voter": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd",
    "proposal_id": "1",
    "option": "Yes"
  }
]
```
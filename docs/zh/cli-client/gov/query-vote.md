# iriscli gov query-vote

## 描述

查询指定提议、指定投票者的投票情况

## 使用方式

```
iriscli gov query-vote [flags]
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] tendermint节点的链ID                                                                                                                 | Yes      |
| --height        |                            | [int] 查询的区块高度                                                                                  |          |
| --help, -h      |                            | 查询命令帮助                                                                                                                                  |          |
| --indent        |                            | 在JSON响应中添加缩进                                                                                                                          |          |
| --ledger        |                            | 使用连接的硬件记账设备                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                                                                                  |          |
| --proposal-id   |                            | [string] 提议ID                                                                                                        | Yes      |
| --trust-node    | true                       | 关闭响应结果校验                                                                                                                    |          |
| --voter         |                            | [string] bech32编码的投票人地址                                                                                                                        | Yes      |

## 例子

### 查询投票

```shell
iriscli gov query-vote --chain-id=test --proposal-id=1 --voter=faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd
```

通过指定提议、指定投票者查询投票情况。

```txt
{
  "voter": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd",
  "proposal_id": "1",
  "option": "Yes"
}
```
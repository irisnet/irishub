# iriscli upgrade query-switch

## 描述

查询switch信息来知道某人对某升级提议是否发送了switch消息。

## 用法

```
iriscli upgrade query-switch --proposal-id <proposalID> --voter <voter address>
```

## 标志

| 名称, 速记       | 默认值                     | 描述                                                        | 必需     |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |
| --proposal-id      |        | 软件升级提议的ID                              | 是     |
| --voter     |                            | 签名switch消息的地址                             | 是      |
| --chain-id      |                            | [string] tendermint节点的链ID                               | 是       |
| --height        | 最近可证明区块高度           | [int] 查询的区块高度                                              |          |
| --help, -h      |                            | 查询命令帮助                                                |          |
| --indent        |                            | 在JSON格式的应答中添加缩进                                   |          |
| --ledger        |                            | 使用连接的硬件记账设备                                       |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口> |          |
| --trust-node    | true                       | 关闭响应结果校验                                            |          |

## 例子

查询用户`faa1qvt2r6hh9vyg3kh4tnwgx8wh0kpa7q2lsk03fe`是否对升级提议（ID为5）发送了switch消息

```
iriscli upgrade query-switch --proposal-id=5 --voter=faa1qvt2r6hh9vyg3kh4tnwgx8wh0kpa7q2lsk03fe
```

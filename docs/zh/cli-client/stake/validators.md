# iriscli stake validators

## 描述

查询所有验证者

## 用法

```
iriscli stake validators [flags]
```

## 标志

| 名称, 速记       | 默认值                     | 描述                                                                 | 必需     |
| --------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Tendermint节点的链ID                                        |          |
| --height        | 最新的可证明区块高度         | 查询的区块高度                                                       |          |
| --help, -h      |                            | validators命令帮助                                                   |          |
| --indent        |                            | 在JSON响应中添加缩进                                                  |          |
| --ledger        |                            | 使用连接的硬件记账设备                                                 |          |
| --node          | tcp://localhost:26657      | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                    |          |
| --trust-node    | true                       | 关闭响应结果校验                                                       |          |

## E例子

### 查询验证者

```shell
iriscli stake validators
```

运行成功以后，返回的结果如下：

```txt
Validator
Operator Address: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Validator Consensus Pubkey: fvp1zcjduepq47906n2zvq597vwyqdc0h35ve0fcl64hwqs9xw5fg67zj4g658aqyuhepj
Jailed: false
Status: Bonded
Tokens: 100.0000000000
Delegator Shares: 100.0000000000
Description: {node0   }
Bond Height: 0
Unbonding Height: 0
Minimum Unbonding Time: 1970-01-01 00:00:00 +0000 UTC
Commission: {{0.0000000000 0.0000000000 0.0000000000 0001-01-01 00:00:00 +0000 UTC}}
```

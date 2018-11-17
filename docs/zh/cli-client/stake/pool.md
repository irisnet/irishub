# iriscli stake pool

## 描述

查询最新的权益池

## 用法

```
iriscli stake pool [flags]
```

## 标志

| 名称, 速记                  | 默认值                     | 描述                                                                | 必需     |
| -------------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id                 |                            | [string] Tendermint节点的链ID                                        |          |
| --height                   | 最新的可证明区块高度         | 查询的区块高度                                                       |          |
| --help, -h                 |                            | pool命令帮助                                                         |          |
| --indent                   |                            | 在JSON响应中添加缩进                                                  |        |
| --ledger                   |                            | 使用连接的硬件记账设备                                                |          |
| --node                     | tcp://localhost:26657      | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                  |          |
| --trust-node               | true                       | 关闭响应结果校验                                                      |          |

## 例子

### 查询最新的权益池

```shell
iriscli stake pool
```

运行成功以后，返回的结果如下：

```txt
Pool
Loose Tokens: 49.8289125612
Bonded Tokens: 100.1800000000
Token Supply: 150.0089125612
Bonded Ratio: 0.6678269863
```

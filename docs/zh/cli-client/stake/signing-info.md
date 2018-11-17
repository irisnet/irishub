# iriscli stake signing-info

## 描述

查询验证者签名信息

## 用法

```
iriscli stake signing-info [validator-pubkey] [flags]
```

## 标志

| 名称, 速记                  | 默认值                     | 描述                                                                 | 必需     |
| -------------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id                 |                            | [string] Tendermint节点的链ID                                        |          |
| --height                   | 最新的可证明区块高度         | 查询的区块高度                                                       |          |
| --help, -h                 |                            | signing-info命令帮助                                                 |          |
| --indent                   |                            | 在JSON响应中添加缩进                                                  |          |
| --ledger                   |                            | 使用连接的硬件记账设备                                                |          |
| --node                     | tcp://localhost:26657      | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                   |          |
| --trust-node               | true                       | 关闭响应结果校验                                                      |          |

## 例子

### 查询验证者签名信息

```shell
iriscli stake signing-info ValidatorPublicKey
```

运行成功以后，返回的结果如下：

```txt
Start height: 0, index offset: 2136, jailed until: 1970-01-01 00:00:00 +0000 UTC, missed blocks counter: 0
```
